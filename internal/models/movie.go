package models

type Movie struct {
	ID			   string
	Title          string
	Description    string
	Voiceover      string
	Subtitles      string
	Quality        string
	ProductionYear uint
	Country        string
	Slogan         string
	Director       string
	Scriptwriter   string
	Producer       string
	Operator       string
	Composer       string
	Artist         string
	Montage        string
	Budget         string
	Duration       string
	Actors         []string
}
