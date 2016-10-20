package dashboard

import "github.com/stretchr/testify/mock"

type DatabaseManagerMock struct {
	mock.Mock
}

func (o *DatabaseManagerMock) GetAll(pageSize, offset int) (result []Widget, err error) {
	args := o.Called(pageSize, offset)
	return args.Get(0).([]Widget), args.Error(1)
}

func (o *DatabaseManagerMock) GetPages() (result []Page, err error) {
	args := o.Called()
	return args.Get(0).([]Page), args.Error(1)
}

func (o *DatabaseManagerMock) InsertWidget(widget *Widget) (int64, error) {
	args := o.Called(widget)
	return args.Get(0).(int64), args.Error(1)
}

func (o *DatabaseManagerMock) InsertOrUpdateWidget(widget *Widget) (int64, error) {
	args := o.Called(widget)
	return args.Get(0).(int64), args.Error(1)
}

func (o *DatabaseManagerMock) InsertWidgetToPage(pageId int64, widgetId string) (int64, error) {
	args := o.Called(pageId, widgetId)
	return args.Get(0).(int64), args.Error(1)
}

func (o *DatabaseManagerMock) InsertPage(page *Page) (int64, error) {
	args := o.Called(page)
	return args.Get(0).(int64), args.Error(1)
}

func (o *DatabaseManagerMock) UpdatePage(page *Page) (int64, error) {
	args := o.Called(page)
	return args.Get(0).(int64), args.Error(1)
}

func (o *DatabaseManagerMock) GetPageWidgets(pageId int64) (result []Widget, err error) {
	args := o.Called(pageId)
	return args.Get(0).([]Widget), args.Error(1)
}

func (o *DatabaseManagerMock) GetUnlinkedWidgets() (result []Widget, err error) {
	args := o.Called()
	return args.Get(0).([]Widget), args.Error(1)
}

func (o *DatabaseManagerMock) Close() error {
	args := o.Called()
	return args.Error(0)
}
