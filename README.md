webhook-validation-proxy
===============================================
Golang validation proxy for [github webhooks][github-webhooks]. It checks headers and payload-secret signature and if all is correct request is proxied to target.

[![Build Status](https://travis-ci.org/mslusarczyk/webhook-validation-proxy.svg?branch=master)](https://travis-ci.org/mslusarczyk/webhook-validation-proxy)

Setup
-----------------------------------------------
Configuration params can be defined as command line args, env variables and via config file. 

For more details see [flag][flag].

```
Usage of webhook-validation-proxy:
  -config string
    	Path to config file
  -context string
    	Context path for proxy for webhooks handling (default "/github-webhook/")
  -port string
    	Port for proxy to listen on (default "8888")
  -secret string
    	Secret assosiated with GH webhook
  -target string
    	Target address with port (default "http://localhost:8080")
```

[github-webhooks]: https://developer.github.com/webhooks/
[flag]: https://github.com/namsral/flag
