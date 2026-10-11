package model

type IssueReport struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Contact     string `json:"contact"`
}
