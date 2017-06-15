# ceph-dashboard - a free ceph dashboard 
Based on https://github.com/Crapworks/ceph-dash and rewrite using golang.

# Easy deployment
All js, css and image resource files are compiled into a single executive file, it is very easy to deploy.

```bash
./ceph-dashboard  -c /some/path/to/ceph.conf  -http "0.0.0.0:8080"
```
and ceph-dashboard will listen on 8080 port and wait for your request. 
Just open http://address:8080/ in your browser.

# Screenshot
![ceph dashboard screenshot](/screenshot/ceph_dashboard_screenshot.png)
