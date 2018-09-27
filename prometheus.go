package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/kataras/iris"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	incomingTotalRequest = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "incoming_request",
		Help: "Total incomming request",
	})
	requestByConfigId = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_by_config_id",
		Help: "Request by config id",
	}, []string{"config_id"})
	failByLabel = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "fail_by_label",
		Help: "Fail description by label",
	}, []string{"label"})
	reactkResponseCode = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "reactk_response_code",
		Help: "Reactk response code",
	}, []string{"code"})
)

func init() {
	prometheus.MustRegister(incomingTotalRequest, requestByConfigId, reactkResponseCode, failByLabel)
	go func() {
		app := iris.New()
		app.Get("prometheus", iris.FromStd(promhttp.Handler()))
		app.Run(iris.Addr(":9191"), iris.WithoutVersionChecker)
	}()

}
