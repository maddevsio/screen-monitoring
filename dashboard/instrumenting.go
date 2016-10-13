package dashboard

import (
	"time"

	gometrics "github.com/rcrowley/go-metrics"
)

type instrumentingMiddleware struct {
	requestCount   gometrics.Counter
	requestLatency gometrics.Histogram
	countResult    gometrics.Histogram
	next           DashboardService
}

func NewInstrumentingMiddleware(requestCount gometrics.Counter,
	requestLatency gometrics.Histogram,
	countResult gometrics.Histogram,
	next DashboardService) DashboardService {
	return instrumentingMiddleware{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		countResult:    countResult,
		next:           next,
	}
}

func (mw instrumentingMiddleware) GetPages() (pages []Page, err error) {
	defer func(begin time.Time) {
		mw.requestCount.Inc(1)
		mw.requestLatency.Update(time.Since(begin).Nanoseconds())
	}(time.Now())

	pages, err = mw.next.GetPages()
	return
}

func (mw instrumentingMiddleware) Register(widget Widget) (pr RegisterResponse, err error) {
	defer func(begin time.Time) {
		mw.requestCount.Inc(1)
		mw.requestLatency.Update(time.Since(begin).Nanoseconds())
	}(time.Now())

	pr, err = mw.next.Register(widget)
	return
}

func (mw instrumentingMiddleware) Init() (errs []error, ok bool) {
	defer func(begin time.Time) {
		mw.requestCount.Inc(1)
		mw.requestLatency.Update(time.Since(begin).Nanoseconds())
	}(time.Now())

	errs, ok = mw.next.Init()
	return
}
