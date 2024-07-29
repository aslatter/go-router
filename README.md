# Go-Router

`go-router` provides a basic wrapper around the standrd-library HTTP-server-mux to allow:

1. Registering of middleware
2. Sub-routers rooted at specific paths, with scoped middleware

The standard-library server-mux is entirely re-used at runtime - this library only
eases the construction of the mux.
