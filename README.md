# kong-go-plugins-demo

#### KONG API Gateway
The world's most popular open source API gateway, built for multi-cloud and hybrid, and optimized for microservices and distributed architectures. It is built on top of a lightweight proxy to deliver unparalleled latency, performance and scalability for microservice applications regardless of where they run.

One of the coolest features Kong offers, it is the possibility of writing your own custom plugins using Golang in order to interact with incoming request and outcoming responses from the plugin's code itself.
Hence, this has been the motivation for this little project to show off how to integrate Golang plugins and Kong API Gateway.

I have been looking through the Kong's documentation, and there is a section for this feature, but it is not clear enough to understand the full picture. Also, I looked up in many places and found just a few examples, but many of them assumes that you have already set some part of the process to achieve this.
So that I have built this project for you to have a reference as a quick start on this topic, everything packed up in docker images so, you just need to run 2 commands and see how this works.

If you follow along the [Makefile](https://github.com/lucas-dev-it/kong-go-plugins-demo/blob/master/Makefile), you will see a pretty  much descriptive set of commands to get everything set up for you to play around with Kong.

#### Dependencies
The only thing you have to care about, is the fact I don't find `cURL` commands quite easy to read, so I used another library to make HTTP calls from the terminal and scripts. I used [httpie](https://github.com/jakubroztocil/httpie) cli. So install it for your OS version, and that would be all set to run all this bundled stuff.

#### Run it
To start up all the containers just do:
```
$ make prepare-and-run
```
Once everything is started you will see:
```
CONTAINER ID        IMAGE                                        COMMAND                  CREATED             STATUS              PORTS                                                                                                NAMES
badddbd57308        lucas-dev-it/kong-go-plugins-demo:latest     "./app"                  21 minutes ago      Up 21 minutes       0.0.0.0:3333->3333/tcp                                                                               login-api-demo
1f439fc70cdb        pantsel/konga                                "/app/start.sh"          21 minutes ago      Up 21 minutes       0.0.0.0:1337->1337/tcp                                                                               konga
5094cebbb171        kong-go-plugins-demo/custom-plugins:latest   "/docker-entrypoint.…"   21 minutes ago      Up 21 minutes       0.0.0.0:8000->8000/tcp, 127.0.0.1:8001->8001/tcp, 0.0.0.0:8443->8443/tcp, 127.0.0.1:8444->8444/tcp   kong
2f14ee1302c9        postgres:9.6                                 "docker-entrypoint.s…"   21 minutes ago      Up 21 minutes       0.0.0.0:5432->5432/tcp                                                                               kong-database
``` 
(*) Notice there a `konga` container running... yes I packed up the UI as well, you can visit it at [http://localhost:1337](http://localhost:1337) and create your own admin account, after that when it asks for the Kong Admin URL just type `http://kong:8001` and that's pretty much it to get it connected to Kong Admin. 

In order to enable the `example` plugin you need to configure your service and route first, and then set the plugin. You can achieve this by typing:
```
$ make setup-endpoints
```
This will run [./setup_endpoints.sh](https://github.com/lucas-dev-it/kong-go-plugins-demo/blob/master/setup_endpoints.sh) script. Which contains all the HTTP calls to the Kong Admin API to set them up.

The last response you see here it is just a mocked response from a tiny mocked server I have included in this project under [_demo](https://github.com/lucas-dev-it/kong-go-plugins-demo/tree/master/_demo/login-api-demo) folder, which will be built when you run the first `make` command.
I tried to make it silly enough just for demo purposes, so the plugin what it does it is just adding a custom header with the path we just hit through the Kong API Gateway.
The custom header is the `x-go-example-path: path /api/users/login` which sits on the headers sections of the below response:
 
```
HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 76
Content-Type: application/json
Date: Sat, 11 Jul 2020 16:05:04 GMT
Via: kong/2.0.1
X-Kong-Proxy-Latency: 37
X-Kong-Upstream-Latency: 2
x-go-example-path: path /api/users/login

{
    "data": {
        "scopes": [
            "inventory",
            "payment",
            "orders",
            "other"
        ]
    },
    "success": true
}

``` 

I will be iteration this example in upcoming weeks just to get something more complex going on with custom plugins.

Thanks, and enjoy it!
