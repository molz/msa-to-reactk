package main

import (
	"flag"
	"github.com/kataras/iris"
	"time"
	"net/http"
	"fmt"
)

var (
	bind              = flag.String("bind", ":8080", "Http listen port :xxxx")
	concurrentRequest = flag.Int("concurrent", 10, "Concurrent requests")
	reactkApi         = flag.String("react-api-url", "", "")
)

func main() {
	flag.Parse()

	app := iris.New()
	app.Post("/api/v1/event", handlerPostEvent)
	app.Run(iris.Addr(*bind), iris.WithoutVersionChecker)
}

func getHttpClient(maxConn int, timeout time.Duration) *http.Client {
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = maxConn
	defaultTransport.MaxIdleConnsPerHost = maxConn
	defaultTransport.DisableKeepAlives = false

	client := http.Client{
		Timeout:   timeout,
		Transport: &defaultTransport,
	}
	return &client
}
