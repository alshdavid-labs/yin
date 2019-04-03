
# Yin

<img align="right" width="200px" src="https://qkgo.github.io/yin/yin.jpg">

### Convenience utilities for idomatic Go HTTP servers

```
go get -u github.com/qkgo/yin
```

This library is compatible with the standard HTTP server in Go,
or any routers that respect it's patterns.
In my examples I am using the Chi router.

Get it here: https://github.com/go-chi/chi

## Getting Started

```Go
package main

import (
    "net/http"
    "github.com/qkgo/yin"
    "github.com/go-chi/chi"
)

func main() {
    r := chi.NewRouter()
    r.Use(yin.SimpleLogger)

    r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
        yin.Res(w, r).JSON(yin.H{
            "message": "world"
        })
    })

    http.ListenAndServe(":3000", r)
}
```


## Logging

Yin's logging format is inspired by Go Gin's logging style.

Choose a preconfigured logger, or configure your own.<br>
Preconfigured routers by default will ignore routes that start with "/ping"

```Go
r.Use(yin.SimpleLogger)
r.Use(yin.DefaultLogger)
r.Use(yin.Logger(os.Stdout, &yin.LoggerConfig{}))
```

## Post Body

Getting body from POST request from map
```Go
r.Post("/incoming-map", func(w http.ResponseWriter, r *http.Request) {
    res, req := yin.Event(w, r)
    body := map[string]interface{}
    req.BindBody(&body)

    fmt.Println(body)

    res.SendStatus(http.StatusNoContent)
})
```

Getting body from POST request from struct
```Go
type request struct {
    Hello string `json:"hello"`
}

r.Post("/incoming-struct", func(w http.ResponseWriter, r *http.Request) {
    res, req := yin.Event(w, r)
    body := request{}
    req.Body(&body)

    fmt.Println(body)

    res.SendStatus(http.StatusNoContent)
})
```

## Headers

Setting response headers

```Go
yin.Res(w, r).
    SetHeader("Key", "Value").
    SendStatus(http.StatusNoContent)
```

There is a convenience struct with common header names

```Go
yin.Res(w, r).
    SetHeader(yin.Headers.Origin, "*").
    SendStatus(http.StatusNoContent)
```

## Serving your Client

Yin does not provide any templating capabilities, it simply allows you to configure
your server to serve a client application.

### Single Page Application

```Go
r.Get("/*", yin.ServeClient(yin.ClientConfig{
    Directory:             "./public",
    SinglePageApplication: true,
}))
```

```Go
r.Get("/base*", yin.ServeClient(yin.ClientConfig{
    Directory:             "./public",
    BaseHref:              "base",
    SinglePageApplication: true,
}))
```

### Static Refresh App

```Go
r.Get("/*", yin.ServeClient(yin.ClientConfig{
    Directory:             "./public",
}))
```

```Go
r.Get("/base*", yin.ServeClient(yin.ClientConfig{
    Directory:             "./public",
    BaseHref:              "base",
}))
```

## Testing Route Handlers

The following is a `POST` request that takes a payload which
has two properties `a` and `b`. The route handler takes the properties,
adds them, then responds with the addition of the two.

```Go
// add-numbers_post.go

type AddNumbersRequest struct {
    A int `json:"a"`
    B int `json:"b"`
}

func AddNumbersHandler(w http.ResponseWriter, r *http.Request) {
    res, req := yin.Event(w, r)

    body := AddNumbersRequest{}
    req.Body(&body)

    result := body.A + body.B

    res.JSON(yin.H{
        "result": result,
    })
}
```

You would use the following test to check it performs as expected

```Go
// add-numbers_post_test.go
import "github.com/qkgo/mock_http"

func TestAddNumbersHandler(t *testing.T) {
    // Create the mock HTTP event
    w := mock_http.NewResponseWriter()
    r := mock_http.NewRequest()

    // Add mock data to it
    r.Body = mock_http.ResponseBody(yin.H{
        "a": 1,
        "b": 1,
    })

    // Run handler with mocks
    AddNumbersHandler(w, r)

    // Get results
    result := w.GetBodyJSON()

    // Assert against results
    if result["result"] == nil && result != 2 {
        t.Errorf("Didn't get the right response")
    }
}
```
