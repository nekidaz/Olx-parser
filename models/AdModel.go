package models

type AdModel struct {
	Title     string `json:"title"`
	Price     string `json:"price"`
	Location  string `json:"location"`
	Condition string `json:"condition"`
	Link      string `json:"link"`
}
