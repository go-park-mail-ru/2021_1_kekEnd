package models

type Movie struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Voiceover      string `json:"voiceover"`
	Subtitles      string `json:"subtitles"`
	Quality        string `json:"quality"`
	ProductionYear uint `json:"production_year"`
	Country        string `json:"country"`
	Slogan         string `json:"slogan"`
	Director       string `json:"director"`
	Scriptwriter   string `json:"scriptwriter"`
	Producer       string `json:"producer"`
	Operator       string `json:"operator"`
	Composer       string `json:"composer"`
	Artist         string `json:"artist"`
	Montage        string `json:"montage"`
	Budget         string `json:"budget"`
	Duration       string `json:"duration"`
	Actors         []string `json:"actors"`
}
