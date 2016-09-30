package dashboard

import "encoding/json"

type PageContent struct {
	Widgets []Widget `json:"widgets"`
}

type Page struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Visible bool   `json:"visible"`
}

func (w *Page) ToJson() ([]byte, error) {
	return json.Marshal(w)
}

type Widget struct {
	Id      string `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Content string `json:"content"`
	Url     string `json:"url"`
}

func (w *Widget) ToJson() ([]byte, error) {
	return json.Marshal(w)
}
