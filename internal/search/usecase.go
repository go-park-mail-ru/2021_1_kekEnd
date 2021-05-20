package search

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"

type UseCase interface {
	Search(query string) (models.SearchResult, error)
}