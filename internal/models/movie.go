package models

type Movie struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Voiceover      []string `json:"voiceover"`
	Subtitles      []string `json:"subtitles"`
	Quality        string   `json:"quality"`
	ProductionYear uint     `json:"production_year"`
	Country        []string `json:"country"`
	Genre          []string `json:"genre"`
	Slogan         string   `json:"slogan"`
	Director       string   `json:"director"`
	Scriptwriter   string   `json:"scriptwriter"`
	Producer       string   `json:"producer"`
	Operator       string   `json:"operator"`
	Composer       string   `json:"composer"`
	Artist         string   `json:"artist"`
	Montage        string   `json:"montage"`
	Budget         string   `json:"budget"`
	Duration       string   `json:"duration"`
	Actors         []string `json:"actors"`

	Poster         string `json:"poster"`
	Banner         string `json:"banner"`
	TrailerPreview string `json:"trailer_preview"`
}
