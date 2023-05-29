package models

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"strings"
	"time"
)

type Metric struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
	Buckets         []float64
	Objectives      map[float64]float64
}

var reqCnt = &Metric{
	ID:          "reqCnt",
	Name:        "requests_total",
	Description: "the number of HTTP requests processed",
	Type:        "counter_vec",
	Args:        []string{"method", "status", "url"}}

var reqSum = &Metric{
	ID:          "reqSum",
	Name:        "request_processing_time_summary_ms",
	Description: "request_processing_time_summary_ms of HTTP requests processed",
	Type:        "summary_vec",
	Args:        []string{"method", "status", "url"},
	Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.09, 0.99: 0.099},
}

var dbSum = &Metric{
	ID:          "dbSum",
	Name:        "database_processing_time_summary_ms",
	Description: "database_processing_time_summary_ms of HTTP requests processed",
	Type:        "summary_vec",
	Args:        []string{"operation", "table", "target"},
	Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.09, 0.99: 0.099},
}

var reqHis = &Metric{
	ID:          "reqHis",
	Name:        "request_processing_time_histogram_ms",
	Description: "request_processing_time_histogram_ms of HTTP requests processed",
	Buckets:     prometheus.LinearBuckets(0, 50, 30),
	Args:        []string{"method", "status", "url"},
}

var Prom *Prometheus

type Prometheus struct {
	reqCnt        *prometheus.CounterVec
	timeSummary   *prometheus.SummaryVec
	DbTimeSummary *prometheus.SummaryVec
	timeHistogram *prometheus.HistogramVec
	router        *gin.Engine
	Metric        *Metric
	MetricsPath   string
}

func NewPrometheus(subsystem string) *Prometheus {
	p := &Prometheus{
		Metric:      reqCnt,
		MetricsPath: "/metrics",
	}
	p.registerMetrics(subsystem)
	return p
}

func (p *Prometheus) registerMetrics(subsystem string) {
	metric := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      reqCnt.Name,
			Help:      reqCnt.Description,
		},
		reqCnt.Args,
	)
	p.reqCnt = metric
	reqCnt.MetricCollector = metric

	timeSum := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       reqSum.Name,
			Objectives: reqSum.Objectives,
		},
		reqSum.Args,
	)
	p.timeSummary = timeSum
	reqSum.MetricCollector = timeSum
	dbTimeSum := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       dbSum.Name,
			Objectives: dbSum.Objectives,
		},
		dbSum.Args,
	)
	p.DbTimeSummary = dbTimeSum
	dbSum.MetricCollector = dbTimeSum

	timeHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    reqHis.Name,
			Buckets: reqHis.Buckets,
		},
		reqHis.Args,
	)
	p.timeHistogram = timeHistogram
	reqHis.MetricCollector = timeHistogram

	prometheus.MustRegister(metric)
	prometheus.MustRegister(timeSum)
	prometheus.MustRegister(dbTimeSum)
	prometheus.MustRegister(timeHistogram)
	//prometheus.Unregister(collectors.NewGoCollector())
}

func (p *Prometheus) Use(e *gin.Engine) {
	e.Use(p.handlerFunc())
	p.router = e
	p.setMetricsPath(e)
}
func (p *Prometheus) handlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		if c.Request.URL.String() == p.MetricsPath {
			c.Next()
			return
		}
		c.Next()
		duration := time.Since(start)
		status := strconv.Itoa(c.Writer.Status())
		var url string
		if urlParts := strings.Split(strings.Split(c.Request.URL.String(), "?")[0], "/"); len(urlParts) > 1 {
			url = "/" + urlParts[1]
			if len(urlParts) > 2 {
				url = url + "/" + urlParts[2]
			}
		} else {
			url = c.Request.URL.String()
		}

		p.reqCnt.WithLabelValues(c.Request.Method, status, url).Inc()
		p.timeSummary.WithLabelValues(c.Request.Method, status, url).Observe(float64(duration.Milliseconds()))
		p.timeHistogram.WithLabelValues(c.Request.Method, status, url).Observe(float64(duration.Milliseconds()))
	}
}

func (p *Prometheus) setMetricsPath(e *gin.Engine) {
	p.router.GET(p.MetricsPath, prometheusHandler())
	//	go p.router.Run(p.listenAddress)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
