package dashboard

import (
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

func (s *DashboardServiceTestSuite) TearDownTest() {
	s.MigratorInstance = nil
	s.DbManagerInstance = nil
}

func (suite *DashboardServiceTestSuite) TearDownSuite() {
}
