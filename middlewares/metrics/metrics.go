package metrics

import (
	"time"

	mt "github.com/najeira/goutils/metrics"
	"github.com/najeira/hikaru"
)

var (
	metricsHttp *mt.MetricsHttp
)

func init() {
	metricsHttp = mt.NewMetricsHttp()
	go metricsUpdateClients()
}

func HandlerFunc(h hikaru.HandlerFunc) hikaru.HandlerFunc {
	return func(c *hikaru.Context) {
		metricsStart()
		defer metricsEnd(c, time.Now())

		// call handler
		h(c)
	}
}

func metricsUpdateClients() {
	for range time.Tick(time.Second) {
		metricsHttp.UpdateClients()
	}
}

func metricsStart() {
	metricsHttp.IncClient()
}

func metricsEnd(c *hikaru.Context, start time.Time) {
	metricsHttp.DecClient()
	metricsHttp.Measure(time.Now().Sub(start))

	sc := c.Status()
	if 200 <= sc && sc <= 299 {
		metricsHttp.Mark2xx()
	} else if 300 <= sc && sc <= 399 {
		metricsHttp.Mark3xx()
	} else if 400 <= sc && sc <= 499 {
		metricsHttp.Mark4xx()
	} else if 500 <= sc && sc <= 599 {
		metricsHttp.Mark5xx()
	}
}

func Metrics() map[string]float64 {
	return metricsHttp.Get()
}
