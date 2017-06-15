// Copyright (c) 2017. All rights reserved.
// Author: yaozongyou@vip.qq.com

package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ceph/go-ceph/rados"
)

const version = "0.0.1"

var (
	clusterName = flag.String("clusterName", "Ceph", "cluster name")
	httpListen  = flag.String("http", "0.0.0.0:8999", "host:port to listen on")
	prefix      = flag.String("prefix", "", "prefix")
)

var showVersion bool
var cephConf string
var layoutTmpl *template.Template
var indexTmpl *template.Template

func init() {
	flag.StringVar(&cephConf, "cephconf", "/etc/ceph/ceph.conf", "ceph conf file location")
	flag.StringVar(&cephConf, "c", "/etc/ceph/ceph.conf", "ceph conf file location")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")

	loadBundle()

	layoutTmpl = template.New("").Funcs(template.FuncMap{"HasField": HasField, "SprintField": SprintField})
	template.Must(layoutTmpl.Parse(_bundle["templates/layout.tmpl"]))

	indexTmpl = template.Must(layoutTmpl.Clone())
	template.Must(indexTmpl.Parse(_bundle["templates/index.tmpl"]))
}

type Location struct {
	Name string
	Href string
}

func main() {
	flag.Parse()
	if showVersion {
		fmt.Println(version)
		return
	}

	http.HandleFunc(*prefix+"/", rootHandler)
	http.HandleFunc(*prefix+"/static/", bundleStatic)
	http.HandleFunc(*prefix+"/api/cluster_status/", clusterStatusHandler)

	log.Fatal(http.ListenAndServe(*httpListen, nil))
}

func loadBundle() {
	for k, v := range _bundle {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			fmt.Println(err)
		} else {
			_bundle[k] = string(b)
		}
	}
}

func bundleStatic(w http.ResponseWriter, req *http.Request) {
	mime := mime.TypeByExtension(filepath.Ext(req.URL.Path))
	if mime != "" {
		w.Header().Set("Content-Type", mime)
	}

	f, ok := _bundle[strings.TrimLeft(req.URL.Path, *prefix+"/")]
	if ok {
		w.Write([]byte(f))
	} else {
		http.NotFound(w, req)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	err := indexTmpl.Execute(w, struct {
		ClusterName string
		Prefix      string
		Locations   []Location
		TotalActive string
	}{
		ClusterName: *clusterName,
		Prefix:      *prefix,
		Locations:   []Location{{"Status", *prefix + "/"}},
		TotalActive: "active",
	})

	if err != nil {
		log.Println(err)
	}
}

func clusterStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	clusterStatus, err := getClusterStatus()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get cluster status: %v", err), 500)
		return
	}

	w.Write(clusterStatus)
}

type Entry struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Host   string `json:"host"`
}

func findHostForOsd(osdId float64, osdTree interface{}) string {
	m := osdTree.(map[string]interface{})

	for _, v := range m["nodes"].([]interface{}) {
		node := v.(map[string]interface{})
		osdType := node["type"].(string)

		if osdType == "host" {
			children := node["children"].([]interface{})
			for _, child := range children {
				if osdId == child.(float64) {
					return node["name"].(string)
				}
			}
		}
	}

	return "unknown"
}

func getUnhealthyOsdDetails(osdTree interface{}) (unhealthyOsds []Entry) {
	m := osdTree.(map[string]interface{})

	for _, v := range m["nodes"].([]interface{}) {
		node := v.(map[string]interface{})
		osdType := node["type"].(string)

		if osdType == "osd" {
			if node["exists"].(float64) == 0 {
				continue
			}
			if node["status"].(string) == "down" || node["reweight"].(float64) == 0.0 {
				unhealthyOsds = append(unhealthyOsds, Entry{
					Name:   node["name"].(string),
					Status: node["status"].(string),
					Host:   findHostForOsd(node["id"].(float64), osdTree),
				})
			}
		}
	}

	return
}

func getClusterStatus() ([]byte, error) {
	conn, _ := rados.NewConn()
	conn.ReadConfigFile(cephConf)
	err := conn.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Shutdown()

	resp, _, err := conn.MonCommand([]byte(`{"prefix": "status", "format": "json-pretty"}`))
	if err != nil {
		return nil, err
	}

	var clusterStatus interface{}
	err = json.Unmarshal(resp, &clusterStatus)
	if err != nil {
		return nil, err
	}

	m := clusterStatus.(map[string]interface{})
	osdMap := m["osdmap"]
	m = osdMap.(map[string]interface{})
	osdMap = m["osdmap"]
	m = osdMap.(map[string]interface{})
	numOsds := m["num_osds"].(float64)
	numUpOsds := m["num_up_osds"].(float64)
	numInOsds := m["num_in_osds"].(float64)

	if numUpOsds < numOsds || numInOsds < numOsds {
		resp, _, err := conn.MonCommand([]byte(`{"prefix": "osd tree", "format": "json-pretty"}`))
		if err != nil {
			return nil, err
		}

		var osdTree interface{}
		err = json.Unmarshal(resp, &osdTree)
		if err != nil {
			return nil, err
		}

		var arr []Entry
		unhealthyOsdDetails := getUnhealthyOsdDetails(osdTree)
		for _, v := range unhealthyOsdDetails {
			arr = append(arr, v)
		}

		(clusterStatus.(map[string]interface{}))["osdmap"].(map[string]interface{})["details"] = arr
	}

	return json.MarshalIndent(clusterStatus, "", "\t")
}
