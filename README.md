servedir
========

A very small httpd.

```
Usage of servedir:
  -addr=":8080": [IP Address and] port on which to listen.
  -dir="./": Directory to serve.
  -ip4=false: Use IPv4 (default).
  -ip6=false: Use IPv6.
```
Examples:
```
./servedir -addr=127.0.0.1:4444 -dir=./xfer
./servedir -dir=./
```
