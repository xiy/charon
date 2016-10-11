# Charon

## About

Charon is designed to act as an API Gateway for a service-oriented architecture. It's my first personal attempt at a useful piece of software written in Go. Because of that, the code might not be great and there will be little to no testing.

## Architecture

Charon acts as both a router and a proxy. Given a set of defined services that each have a defined routing rule, it will proxy any HTTP requests that match a rule to the respective service.

## Config

Config is handled entirely through a TOML-based config file. It's very simple and has only a few options:

```tool
port = "9000"
service_timeout = "30s"

[services]
  [services.member]
  url = "http://users.charon.io"
  prefix = "/users/*"
  
```

The value for `prefix` is important, as it defines how the internal router will route requests to each specific service. With the above config, any request to `http://gateway.charon.io` with a path that matches `/users`, e.g. `http://gateway.charon.io/users/1` will be proxied to the service at `http://users.charon.io/users/1`.