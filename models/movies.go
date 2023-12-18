package models

type Movies struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Actors  string `json:"actors"`
	Rating  int    `json:"rating"`
	Details string `json:"details"`
	Genre   string `json:"genre"`
	Saison  int    `json:"saison"`
	Episode int    `json:"episode"`
}
