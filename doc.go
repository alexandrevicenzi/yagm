// Package yagm implements a simple regular expression pattern match mux.
//
// Regular Expressions patterns
//
// YagmMux uses regular expressions to match an URL.
// If you're familiar with Django, the patterns match is almost the same.
// If you're not, take a look in the pattern match explanation below.
//
// Remember, the declared order can affect the match. For example, if you place "^/[a-z]+"
// before "^/index" it will never match "^/index", as the regular expression of "^/[a-z]+"
// can match "/index" too.
//
// The route "^/$" will match only
//   /
//
// The route "^/files" will match
//    /files
//    /files/
//    /files/myfilename.txt
//
// The route "^/index$" will match
//    /index
// and will not match
//    /index/
//    /index/other
//
// If you use named groups like "^/(?P<name>[a-zA-Z]+)$"
// you can use yagm.Param or yagm.Params function to retrieve the value ot the groups.
//
// Example
//
// YagmMux is very close to http.ServeMux.
// You can replace your actual http.ServeMux by YagmMux without any problem.
//
//    package main
//    
//    import (
//        "fmt"
//        "net/http"
//    
//        "github.com/alexandrevicenzi/yagm"
//    )
//    
//    func homeHandler(w http.ResponseWriter, r *http.Request) {
//        fmt.Fprint(w, "Hello World!")
//    }
//    
//    func helloHandler(w http.ResponseWriter, r *http.Request) {
//        name, ok := yagm.Param(r, "name")
//    
//        if !ok {
//            name = "Unknown"
//        }
//    
//        fmt.Fprintf(w, "Hello %s!", name)
//    }
//    
//    func main() {
//        mux := yagm.New()
//        mux.HandleFunc("^/$", homeHandler)
//        mux.Handle("^/files", http.StripPrefix("/files", http.FileServer(http.Dir("./"))))
//        mux.HandleFunc("^/(?P<name>[a-zA-Z]+)$", helloHandler)
//        http.ListenAndServe(":8000", mux)
//    }
//
package yagm
