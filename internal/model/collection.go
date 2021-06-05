package model

type Collection struct {
	Total     int         `json:"total"`
	FromIndex int         `json:"from"`
	ToIndex   int         `json:"to"`
	PerPage   int         `json:"per_page"`
	LastPage  int         `json:"last_page"`
	Data      interface{} `json:"data"`
}
