package dashboard

import "sync"

type RegisterResponse struct {
	Success bool
}

type InsertPageResponse struct {
	Id      int64
	Success bool
}

type DashboardService interface {
	GetPages() (pc []Page, err error)
	Register(widget Widget) (pr RegisterResponse, err error)
	RegisterToPage(pageId int64, widgetId string) (pr RegisterResponse, err error)
	InsertPage(page Page) (response InsertPageResponse, err error)
	Init() ([]error, bool)
}

type dashboardService struct {
	sync.RWMutex
	migrator  Migrator
	dbManager DatabaseManager
}

func NewDashboardService(migrator Migrator, dbManager DatabaseManager) DashboardService {
	return &dashboardService{
		migrator:  migrator,
		dbManager: dbManager,
	}
}

func (d *dashboardService) Init() ([]error, bool) {
	return d.migrator.Up()
}

func (d *dashboardService) GetPages() (pages []Page, err error) {
	d.Lock()
	defer d.Unlock()
	pages, err = d.dbManager.GetPages()
	return
}

func (d *dashboardService) RegisterToPage(pageId int64, widgetId string) (pr RegisterResponse, err error) {
	d.Lock()
	defer d.Unlock()

	_, err = d.dbManager.InsertWidgetToPage(pageId, widgetId)
	if err != nil {
		pr = RegisterResponse{Success: false}
		return
	}

	pr = RegisterResponse{Success: true}
	return
}

func (d *dashboardService) Register(widget Widget) (pr RegisterResponse, err error) {
	d.Lock()
	defer d.Unlock()
	_, err = d.dbManager.InsertOrUpdateWidget(&widget)
	if err != nil {
		pr = RegisterResponse{Success: false}
		return
	}

	pr = RegisterResponse{Success: true}
	return
}

func (d *dashboardService) InsertPage(page Page) (response InsertPageResponse, err error) {
	d.Lock()
	defer d.Unlock()
	id, err := d.dbManager.InsertPage(&page)
	if err != nil {
		response = InsertPageResponse{Success: false}
		return
	}

	response = InsertPageResponse{Id: id, Success: true}
	return
}
