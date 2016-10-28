package dashboard

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DashboardServiceTestSuite struct {
	suite.Suite
	MigratorInstance  *MigratorMock
	DbManagerInstance *DatabaseManagerMock
}

func (s *DashboardServiceTestSuite) SetupSuite() {
}

func (s *DashboardServiceTestSuite) SetupTest() {
	s.MigratorInstance = new(MigratorMock)
	s.DbManagerInstance = new(DatabaseManagerMock)
}

func TestDashboardServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DashboardServiceTestSuite))
}

func (s *DashboardServiceTestSuite) TestGetPagesReturnPagesArray() {
	var pages = []Page{
		Page{
			Id:      1,
			Title:   "page1",
			Visible: true,
			Widgets: []Widget{
				Widget{Id: "widget_1", Url: "http://example.com:8080/", Width: 450, Height: 350},
			},
		},
		Page{
			Id:      2,
			Title:   "page2",
			Visible: true,
			Widgets: []Widget{
				Widget{Id: "widget_1", Url: "http://example.com:8081/", Width: 250, Height: 150},
				Widget{Id: "widget_2", Url: "http://example.com:8081/", Width: 250, Height: 150},
			},
		},
	}

	s.DbManagerInstance.On("GetPages").Return(pages, nil)

	var service = NewDashboardService(s.MigratorInstance, s.DbManagerInstance)
	pages, err := service.GetPages()
	s.DbManagerInstance.AssertNumberOfCalls(s.T(), "GetPages", 1)
	assert.NotNil(s.T(), pages)
	assert.Nil(s.T(), err)
}

func (s *DashboardServiceTestSuite) TestRegisterWidgetSuccess() {
	var expectedResponse = RegisterResponse{Success: true}
	var widget = Widget{
		Id:     "widget_1",
		Url:    "http://example.com:8081/",
		Width:  250,
		Height: 150,
	}

	s.DbManagerInstance.On("InsertOrUpdateWidget", &widget).Return(int64(1), nil)
	var service = NewDashboardService(s.MigratorInstance, s.DbManagerInstance)

	actualResponse, err := service.Register(widget)

	s.DbManagerInstance.AssertNumberOfCalls(s.T(), "InsertOrUpdateWidget", 1)
	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Nil(s.T(), err)
}

func (s *DashboardServiceTestSuite) TestRegisterWidgetFail() {
	var expectedResponse = RegisterResponse{Success: false}
	var expectedError = errors.New("Some error")
	var widget = Widget{
		Id:     "widget_2",
		Url:    "http://example.com:8000/",
		Width:  250,
		Height: 150,
	}

	s.DbManagerInstance.On("InsertOrUpdateWidget", &widget).Return(int64(0), expectedError)
	var service = NewDashboardService(s.MigratorInstance, s.DbManagerInstance)

	actualResponse, actualError := service.Register(widget)

	s.DbManagerInstance.AssertNumberOfCalls(s.T(), "InsertOrUpdateWidget", 1)
	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Equal(s.T(), expectedError, actualError)
}

func (s *DashboardServiceTestSuite) TestInsertPageSuccess() {
	expectedResponse := InsertPageResponse{Id: 1, Success: true}
	var page = Page{
		Title:   "Page1",
		Visible: true,
	}

	s.DbManagerInstance.On("InsertPage", &page).Return(int64(1), nil)

	var service = NewDashboardService(s.MigratorInstance, s.DbManagerInstance)
	actualResponse, err := service.InsertPage(page)

	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Nil(s.T(), err)
}

func (s *DashboardServiceTestSuite) TestInsertPageFail() {
	expectedResponse := InsertPageResponse{Success: false}
	expectedError := errors.New("Some error")
	var page = Page{
		Title:   "Page2",
		Visible: true,
	}

	s.DbManagerInstance.On("InsertPage", &page).Return(int64(0), expectedError)

	var service = NewDashboardService(s.MigratorInstance, s.DbManagerInstance)
	actualResponse, err := service.InsertPage(page)

	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Equal(s.T(), expectedError, err)
}

func (s *DashboardServiceTestSuite) TestRegisterWidgetToPageSuccess() {
	expectedResponse := RegisterResponse{Success: true}
	pageId := int64(1)
	widgetName := "widget_1"

	s.DbManagerInstance.On("InsertWidgetToPage", pageId, widgetName).Return(int64(1), nil)

	var service = NewDashboardService(s.MigratorInstance, s.DbManagerInstance)
	actualResponse, err := service.RegisterToPage(pageId, widgetName)

	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Nil(s.T(), err)
}

func (s *DashboardServiceTestSuite) TestRegisterWidgetToPageFail() {
	expectedResponse := RegisterResponse{Success: false}
	expectedError := errors.New("Some error")

	pageId := int64(1)
	widgetName := "widget_1"

	s.DbManagerInstance.On("InsertWidgetToPage", pageId, widgetName).Return(int64(0), expectedError)

	var service = NewDashboardService(s.MigratorInstance, s.DbManagerInstance)
	actualResponse, actualError := service.RegisterToPage(pageId, widgetName)

	assert.Equal(s.T(), expectedResponse, actualResponse)
	assert.Equal(s.T(), expectedError, actualError)
}

func (s *DashboardServiceTestSuite) TestInitUpMigrationsSuccess() {
	expectedErrors := []error{}
	expectedStatus := true
	s.MigratorInstance.On("Up").Return(expectedErrors, expectedStatus)
	var service = NewDashboardService(s.MigratorInstance, s.DbManagerInstance)
	actualErrors, actualStatus := service.Init()
	s.MigratorInstance.AssertNumberOfCalls(s.T(), "Up", 1)
	assert.True(s.T(), actualStatus)
	assert.Equal(s.T(), expectedErrors, actualErrors)
}

func (s *DashboardServiceTestSuite) TestInitUpMigrationsFail() {
	expectedErrors := []error{errors.New("Errors 1"), errors.New("Errors 2")}
	expectedStatus := false
	s.MigratorInstance.On("Up").Return(expectedErrors, expectedStatus)
	var service = NewDashboardService(s.MigratorInstance, s.DbManagerInstance)
	actualErrors, actualStatus := service.Init()
	s.MigratorInstance.AssertNumberOfCalls(s.T(), "Up", 1)
	assert.False(s.T(), actualStatus)
	assert.Equal(s.T(), expectedErrors, actualErrors)
}

func (s *DashboardServiceTestSuite) TestGetUnregisteredWidgetsSuccess() {
	expectedWidgets := []Widget{Widget{Id: "w1", Url: "http://example1.com/", Width: 450, Height: 350}}
	s.DbManagerInstance.On("GetUnlinkedWidgets").Return(expectedWidgets, nil)
	var service = NewDashboardService(s.MigratorInstance, s.DbManagerInstance)
	actualWidgets, actualError := service.GetUnregisteredWidgets()
	s.DbManagerInstance.AssertNumberOfCalls(s.T(), "GetUnlinkedWidgets", 1)
	assert.Nil(s.T(), actualError)
	assert.Equal(s.T(), expectedWidgets, actualWidgets)
}

func (s *DashboardServiceTestSuite) TearDownTest() {
	s.MigratorInstance = nil
	s.DbManagerInstance = nil
}

func (suite *DashboardServiceTestSuite) TearDownSuite() {
}
