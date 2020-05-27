# Redirector

> a simple golang HTTP(S) server that redirects all requests

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Contributing](../CONTRIBUTING.md)

## About <a name = "about"></a>

Why spin up an nginx server on a pod if you can run a standalone application to
return your 301/307 redirects?

## Getting Started <a name = "getting_started"></a>

```shell
# run directly
go run main.go
# or build and run
go build ./... && ./redirector
# docker!
docker run -e URL=https://artero.dev wwmoraes/redirector
```

This server was made with love and KISS principles in mind 🖤

## Usage <a name = "usage"></a>

A single instance will run both HTTP and HTTPS servers in parallel, given that
you provided valid certificate files for TLS. If not, it'll fail to start the
secure one, but will still run the HTTP one.

All configurations are done using environment variables:

| Variable Name | Default Value | Description                                                                                      |
|---------------|---------------|--------------------------------------------------------------------------------------------------|
| URL           |               | destination URL redirect to, e.g. <https://artero.dev>                                           |
| HTTP_HOST     |               | HTTP server host                                                                                 |
| HTTP_PORT     | `8080`        | HTTP server port                                                                                 |
| HTTPS_HOST    |               | HTTPS server host                                                                                |
| HTTPS_PORT    | `8081`        | HTTPS server port                                                                                |
| KEY_FILE      |               | private key path for the URL certificate                                                         |
| KEY           |               | base64-encoded private key contents for the URL certificate                                      |
| CERT_FILE     |               | public key path for the URL certificate (must be the full chain if CA-signed)                    |
| CERT          |               | base64-encoded public key contents for the URL certificate (must be the full chain if CA-signed) |

Empty host equals `0.0.0.0`, i.e. all interfaces will be bound.
