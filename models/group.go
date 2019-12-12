package models

type Group struct {
	Id          int    `json:"id"`
	CreatorId   int
	Title       string `json:"title"`
	Description string `json:"description"`
}

type JsonGroup struct {
	Message string `json:"message"`
	Status string `json:"status"`
	Group Group `json:"tasks"`
}