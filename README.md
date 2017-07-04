# MessageBox
基于golang 和 reids pubsub 实现的sever sent events 服务器，可以用来对接任意平台作为消息推送

## Usage

There is one route to listen to message published:

- `GET /messagebox`

port default: 9977

### 1. need to set four environments 
 ```
REDISADDR=127.0.0.1:6379 // 配置的redis server 地址
REDISDB=0 // 配置的redis DB

```
### 2. html 5 client example

```
        var sse = new EventSource("http://localhost:9977/messagebox");

        sse.onmessage = function(event) {
            console.log(event.data);
            document.getElementById("result").innerHTML+=event.data + "<br>";
        }

        sse.onerror = function(event) {
        console.log(event);
```

## How to dockerize this service

### 1. 为了得到小的image，我们build outside dockerfile

注意编译的系统应选择cross-linux 64，以适应alpine docker image base

```
go build -o MessageBox 

```
### 2. build docker

```
docker build -t lzhao/messagebox

```
### 3. start up service

利用docker-compose启动服务

```
docker-compose up -f docker-compose.yml -d --force-recreate --remove-orphans
```

