# ceph-dashboard - easy deployment ceph dashboard 
[![Release Version](https://img.shields.io/badge/release-1.0.0-red.svg)](https://github.com/yaozongyou/ceph-dashboard/releases)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/yaozongyou/ceph-dashboard/pulls)
Based on [ceph-dash](https://github.com/Crapworks/ceph-dash) and rewrite using golang.

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
