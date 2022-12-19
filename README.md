goweb
----

web framework written in go

## Installation

1.Install go

2.Use below command insatall goweb

```
go get -u github.com/awaketai/goweb
```

## Quick start

buid a execute file

```
cd build
make build
```

then project root file generate a execute file named gw,run it

## Swagger Use

1.run following command generage file releated to swagger 

```
gw swagger gen
```
will be generated some file to app/http/swagger

2.run app server,then serve will be listening port:8080,then browse: localhost:8080/swagger/index.html

