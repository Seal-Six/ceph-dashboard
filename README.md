# ceph-dashboard - easy deployment ceph dashboard 
[![Release Version](https://img.shields.io/badge/release-1.0.0-red.svg)](https://github.com/yaozongyou/ceph-dashboard/releases)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/yaozongyou/ceph-dashboard/pulls)

Based on [ceph-dash](https://github.com/Crapworks/ceph-dash) and rewrite using golang.

# How to build

ceph-dashboard is written using golang, so for building ceph-dashboard, you have to install golang first.

## Install golang

Golang install instructions can be found at https://golang.org/doc/install.

## Install ceph development libraries

The native RADOS library and development headers are expected to be installed.
```bash
yum install librados-devel
yum install librbd-devel
```
### Building ceph-dashboard

```bash
go get github.com/yaozongyou/ceph-dashboard
```

# Easy deployment

All js, css and image resource files are compiled into a single executive file, it is very easy to deploy.
Run ceph dashboard with the following command:

```bash
./ceph-dashboard -c /some/path/to/ceph.conf  -http "0.0.0.0:8080"
```
Ceph-dashboard will listen on 8080 port and wait for your request. 
Just open http://address:8080/ in your browser.

# Screenshot

![ceph dashboard screenshot](/screenshot/ceph_dashboard_screenshot.png)
