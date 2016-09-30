package dashboard

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   DashboardService
}

func NewLoggerService(logger log.Logger, next DashboardService) DashboardService {
	return loggingMiddleware{logger: logger, next: next}
}

func (mw loggingMiddleware) GetPages() (pc PageContent, err error) {
	pc, err = mw.next.GetPages()
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Pages",
			"input", "[No parametgers]",
			"output", fmt.Sprintf("%+v", pc),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return
}

func (mw loggingMiddleware) Register(widget Widget) (pr RegisterResponse, err error) {
	pr, err = mw.next.Register(widget)
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Register",
			"input", fmt.Sprintf("%+v", widget),
			"output", fmt.Sprintf("%+v", pr),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return
}

func (mw loggingMiddleware) Init() (initErrors []error, ok bool) {
	initErrors, ok = mw.next.Init()
	defer func(begin time.Time) {
		for _, err := range initErrors {
			_ = mw.logger.Log(
				"method", "Init",
				"input", "No params",
				"output", ok,
				"err", err,
				"took", time.Since(begin),
			)
		}
	}(time.Now())
	return
}
