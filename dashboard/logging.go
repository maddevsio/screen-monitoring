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

func (mw loggingMiddleware) GetPages() (pages []Page, err error) {
	pages, err = mw.next.GetPages()
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Pages",
			"input", "[No parametgers]",
			"output", fmt.Sprintf("%+v", pages),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return
}

func (mw loggingMiddleware) GetUnregisteredWidgets() (widgets []Widget, err error) {
	widgets, err = mw.next.GetUnregisteredWidgets()
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetUnregisteredWidgets",
			"input", "[No parametgers]",
			"output", fmt.Sprintf("%+v", widgets),
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
func (mw loggingMiddleware) RegisterToPage(pageId int64, widgetId string) (pr RegisterResponse, err error) {
	pr, err = mw.next.RegisterToPage(pageId, widgetId)
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "RegisterToPage",
			"input", fmt.Sprintf("%+v %+v", pageId, widgetId),
			"output", fmt.Sprintf("%+v", pr),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return
}

func (mw loggingMiddleware) InsertPage(page Page) (response InsertPageResponse, err error) {
	response, err = mw.next.InsertPage(page)
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "InsertPage",
			"input", fmt.Sprintf("%+v", page),
			"output", fmt.Sprintf("%+v", response),
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
