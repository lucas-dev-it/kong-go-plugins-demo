# kong-go-plugins-demo

#### KONG API Gateway - LUA and Golang Plugins Demo
The world's most popular open source API gateway, built for multi-cloud and hybrid, and optimized for microservices and distributed architectures. It is built on top of a lightweight proxy to deliver unparalleled latency, performance and scalability for microservice applications regardless of where they run.

One of the coolest features Kong offers, it is the possibility of writing your own custom plugins using Golang in order to interact with incoming request and outcoming responses from the plugin's code itself.
Hence, this has been the motivation for this little project to show off how to integrate Golang plugins and Kong API Gateway.

I have been looking through the Kong's documentation, and there is a section for this feature, but it is not clear enough to understand the full picture. Also, I looked up in many places and found just a few examples, but many of them assumes that you have already set some part of the process to achieve this.
So that I have built this project for you to have a reference as a quick start on this topic, everything packed up in docker images so, you just need to run 2 commands and see how this works.

If you follow along the [Makefile](https://github.com/lucas-dev-it/kong-go-plugins-demo/blob/master/Makefile), you will see a pretty  much descriptive set of commands to get everything set up for you to play around with Kong.
