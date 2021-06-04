package model

type Type struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Group Group  `json:"group"`
}

type Group struct {
	ID       int32    `json:"id"`
	Name     string   `json:"name"`
	Category Category `json:"category"`
}

type Category struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}
