# Go-Router

`go-router` provides a basic wrapper around the standrd-library HTTP-server-mux to allow:

1. Registering of middleware
2. Sub-routers rooted at specific paths, with scoped middleware

The standard-library server-mux is entirely re-used at runtime - this library only
eases the construction of the mux.

## Example

```go
root := router.New()

// if we want the '404' response to go through
// registered middleware we need to handle that
// in the router.
root.Handle("/", http.NotFoundHandler())

// create a sub-router under the '/api' prefix
api := root.New("/api")

// this handler will respond to a GET request at '/api/test'.
api.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
    log.Println("in /api/test handler")
})

// add middleware dedicated to the '/api' prefix
api.Use(func(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("in api middleware")
        h.ServeHTTP(w, r)
    })
})

// add a handler using the 'root' router
root.HandleFunc("GET /test/{id}", func(w http.ResponseWriter, r *http.Request) {
    log.Printf("in /test handler for id: %q", r.PathValue("id"))
})

// add middleware to the 'root' router. This will also apply to requests
// to the '/api' routes.
root.Use(func(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("in root middleware")
        h.ServeHTTP(w, r)
    })
})

// serve up traffic
err := http.ListenAndServe(":8080", root.Handler())
if err != nil {
    fmt.Fprintln(os.Stderr, "error: ", err)
    os.Exit(1)
}
```