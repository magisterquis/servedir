servedir
========

A very small httpd which serves static files.

Yeah, yeah, it's yet another a wrapper around `http.FileServer`.

Please run with `-h` for a complete list of options.

Installation
------------
```sh
$ go get magisterquis/servedir
```
Compiled binaries available upon request.

Examples
--------
Serve via http on port `8080` and https on port `4433` using the keypair from
`cert.pem` and `key.pem`:
```sh
./servedir
```

Serve files only via http on `127.0.0.1:8888` from `/tmp/d`:
```sh
./servedir -https no -http 127.0.0.1:8888 -dir /tmp/d
```

Windows
-------
Works just fine.
