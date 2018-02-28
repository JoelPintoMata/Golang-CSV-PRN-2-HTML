Run
===
```
$ go run main.go
```

Tests
=====
```
$ go test -v
```

Endpoints
=========
```
Verb | URI                         | Example
----------------------------------------------------------------------
GET  | <server>:<port>/convert/csv | http://localhost:8080/convert/csv
GET  | <server>:<port>/convert/prn | http://localhost:8080/convert/prn
```

Installing dependencies
=======================
```
$ go get github.com/gin-gonic/gin
```
```
$ go get gopkg.in/russross/blackfriday.v2
```

Further reading
===============
1. HTTP web framework [gin](https://github.com/gin-gonic/gin)
2. Markdown processor [blackfriday](https://github.com/russross/blackfriday)
