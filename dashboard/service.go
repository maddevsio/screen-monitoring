package dashboard

import (
	"encoding/json"
	"sync"
)

type Widget struct {
	ID      string `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Content string `json:"content"`
	Url     string `json:"url"`
}

func (w *Widget) ToJson() ([]byte, error) {
	return json.Marshal(w)
}

type PageContent struct {
	Widgets []Widget `json:"widgets"`
}

type RegisterResponse struct {
	Success bool
}

type DashboardService interface {
	GetPages() (pc PageContent, err error)
	Register(widget Widget) (pr RegisterResponse, err error)
}

type dashboardService struct {
	registeredWidgets map[string]Widget
	sync.RWMutex
}

func NewDashboardService() DashboardService {
	return &dashboardService{registeredWidgets: make(map[string]Widget, 0)}
}

func (d dashboardService) GetPages() (pc PageContent, err error) {
	result := make([]Widget, 0)

	for _, v := range d.registeredWidgets {
		result = append(result, v)
	}

	pc = PageContent{
		Widgets: result,
	}
	err = nil
	return
}

func (d *dashboardService) Register(widget Widget) (pr RegisterResponse, err error) {
	d.Lock()
	defer d.Unlock()
	d.registeredWidgets[widget.ID] = widget
	pr = RegisterResponse{Success: true}
	err = nil
	return
}
