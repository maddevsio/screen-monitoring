package dashboard

import "github.com/stretchr/testify/mock"

type DashboardServiceMock struct {
	mock.Mock
}

func (o *DashboardServiceMock) GetPages() (pc []Page, err error) {
	args := o.Called()
	return args.Get(0).([]Page), args.Error(1)
}

func (o *DashboardServiceMock) Register(widget Widget) (pr RegisterResponse, err error) {
	args := o.Called(widget)
	return args.Get(0).(RegisterResponse), args.Error(1)
}

func (o *DashboardServiceMock) RegisterToPage(pageId int64, widgetId string) (pr RegisterResponse, err error) {
	args := o.Called(pageId, widgetId)
	return args.Get(0).(RegisterResponse), args.Error(1)
}

func (o *DashboardServiceMock) InsertPage(page Page) (response InsertPageResponse, err error) {
	args := o.Called(page)
	return args.Get(0).(InsertPageResponse), args.Error(1)
}

func (o *DashboardServiceMock) Init() ([]error, bool) {
	args := o.Called()
	return args.Get(0).([]error), args.Bool(1)
}
