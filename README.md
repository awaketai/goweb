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

## Set env

webserver will read root dir .env file,the env contains:development | testing  | production

## App server

### 1.Start app server
app will read command ADDRESS parameter,if this parameter not set,the config dir app.yaml file will be read,and the address parameter will be load,if above parameter not set both,default address:8080

1.env command address

```
ADDRESS=:8081 ./gw app start
```
2. app.yaml address

```
config/development:
address: 8081
```
3.command address

```
./gw app start --address=:8081

```

**start app server daemon**

```
./gw app start --daemon=true
or 
./gw app start --d=true
```

```
./gw app start 
```
### Restart app server

```
./gw app restart
```

### Stop app server

```
./gw app stop
```

## Swagger Use

1.run following command generage file releated to swagger 

```
gw swagger gen
```
will be generated some file to app/http/swagger

2.run app server,then serve will be listening port:8080,then browse: localhost:8080/swagger/index.html

