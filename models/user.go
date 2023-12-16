package models

type User struct {
	Id      int    `json:"id"`
	User    string `json:"user"`
	Pswd    string `json:"pswd"`
	Account string `json:"account"`
	Info    Info
}

type Info struct {
	Name      string `json:"name"`
	Firstname string `json:"firstname"`
	Mail      string `json:"mail"`
	Cell      int    `json:"cell"`
	Adress    string `json:"adress"`
}
