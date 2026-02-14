package main

type Pokemon struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Types interface{} `json:"types"`
}
