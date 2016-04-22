# yagm [![Build Status](https://travis-ci.org/alexandrevicenzi/yagm.svg?branch=master)](https://travis-ci.org/alexandrevicenzi/yagm) [![GoDoc](https://godoc.org/github.com/alexandrevicenzi/yagm?status.svg)](http://godoc.org/github.com/alexandrevicenzi/yagm)

yagm is an acronym of Yet Another Go Mux.

YagmMux uses regular expressions to match an URL.
If you're familiar with Django, the patterns match is almost the same.

YagmMux is also very close to http.ServeMux.
You can replace your actual http.ServeMux by YagmMux without any problem.

## Install

`go get github.com/alexandrevicenzi/yagm`

## Usage

You can register your handlers as in `http.ServeMux`.

```go
mux := yagm.New()
mux.HandleFunc("^/$", homeHandler)
mux.Handle("^/files", http.StripPrefix("/files", http.FileServer(http.Dir("./"))))
mux.HandleFunc("^/(?P<name>[a-zA-Z]+)$", helloHandler)
```

You can get the value of named groups with

```go
params, ok := yagm.Params(request)
value, ok := params["group_name"]
```

or

```go
value, ok := yagm.Param(request, "group_name")
```

## Full Example

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/alexandrevicenzi/yagm"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello World!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    name, ok := yagm.Param(r, "name")

    if !ok {
        name = "Unknown"
    }

    fmt.Fprintf(w, "Hello %s!", name)
}

func main() {
    mux := yagm.New()
    mux.HandleFunc("^/$", homeHandler)
    mux.Handle("^/files", http.StripPrefix("/files", http.FileServer(http.Dir("./"))))
    mux.HandleFunc("^/(?P<name>[a-zA-Z]+)$", helloHandler)
    http.ListenAndServe(":8000", mux)
}
```

## License

MIT
