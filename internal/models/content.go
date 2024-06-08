package models

type Content struct {
	PageTitle string    `json:"page_title"`
	Results   []*Result `json:"endpoints"`
}
