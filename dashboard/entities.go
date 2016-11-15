package dashboard

import "encoding/json"

type PageContent struct {
	Widgets []Widget `json:"widgets"`
}

type Page struct {
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
	Visible bool     `json:"visible"`
	Widgets []Widget `json:"content"`
}

func (w *Page) ToJson() ([]byte, error) {
	return json.Marshal(w)
}

type Widget struct {
	Id      *string `json:"id"`
	Width   *int64  `json:"width"`
	Height  *int64  `json:"height"`
	Content *string `json:"content"`
	Url     *string `json:"url"`
}

func NewWidget(id string, width, height int64, url string) Widget {
	var defaultContent string
	return Widget{
		Id:      &id,
		Width:   &width,
		Height:  &height,
		Url:     &url,
		Content: &defaultContent,
	}
}

func (w *Widget) ToJson() ([]byte, error) {
	return json.Marshal(w)
}
