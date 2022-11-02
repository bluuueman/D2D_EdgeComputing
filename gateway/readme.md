# To build this porject you need to install gin and go4vl
```
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go get github.com/vladimirvivien/go4vl/v4l2
go get -u github.com/gin-gonic/gin
```
# How to use this gateway ?
## 1.Service Registry
### Send your service info to gateway
```
URL             ip:port/service
Method          POST
Content-Type    application/json
Body
{
    "ip":"192.168.0.1",
    "priority":4,
    "data":{
        "1":{
            "service":"s1",
            "port":"8080"
        },
        "2":{
            "service":"s2",
            "prot":"8088"
        }
    }
}
```
## 2.Send Heartbeat
### After registration, send server status as heartbeat every few seconds
```
URL             ip:port/server
Method          POST
Content-Type    application/json
Body
{
    "ip":"192.168.0.1",
    "priority":4
}
```
## 3.Set Job ( For test only )
### Set the job you want to offload to server
```
URL             ip:port/job
Method          POST
Content-Type    application/json
Body
{
    "service":"service name"
}
```