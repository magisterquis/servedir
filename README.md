servedir
========

A very small httpd which serves static files.

This update breaks the previous interface (again).  On the other hand, I can't
imagine anybody actually used the previous code (and it's not likely anybody
will use this, either).

Yeah, yeah, it's a wrapper around `http.FileServer`.

Also, this documentation needs an update.

Installation
------------
```sh
$ go get magisterquis/servedir
```
If you don't have go and desperately need a static file server, I'm happy
to send you compiled binaries.

Examples
--------
Serve files from the current directory which holds the TLS cert and key.  In
practice, this is probably a lousy idea.

```sh
$ servedir
2016/03/26 18:19:49 Serving files from .
2016/03/26 18:19:49 Listening on [::]:4433 (tcp6) for HTTPS requests
2016/03/26 18:19:49 Listening on 0.0.0.0:8080 (tcp4) for HTTP requests
2016/03/26 18:19:49 Listening on [::]:8080 (tcp6) for HTTP requests
2016/03/26 18:19:49 Listening on 0.0.0.0:4433 (tcp4) for HTTPS requests
```

Serve files only on port 80, HTTP, from `/tmp`.  This probably is a pretty
bad idea as well.
```sh
$ servedir -dir /tmp -nohttps -http 80
```

Securely serve files from your personal htdocs directory.
```sh
$ servedir -dir ~/htdocs -nohttp
```

Usage
-----
```
Usage servedir [options]

Serves files from the given directory.  Content-type will be automatically
determined.

Options:
  -4	Don't listen on IPv4
  -6	Don't listen on IPv6
  -cert certificate
    	HTTPS certificate (ignored if -nohttps is given) (default "cert.pem")
  -dir directory
    	Files in this directory will be served (default ".")
  -http address
    	HTTP [address and] port (default "8080")
  -https address
    	HTTPS listen [address and] port (default ":4433")
  -key key
    	HTTPS key (ignored if -nohttps is given) (default "key.pem")
  -nohttp
    	Don't handle HTTP requests
  -nohttps
    	Do not handle HTTPS requests
```

Windows
-------
It should work on Windows just fine, though you'll have to use the goofy
Windows backwards slashes.
