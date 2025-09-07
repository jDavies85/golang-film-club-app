// internal/usecase/search_movies.go
package usecase

import (
	"context"

	"github.com/jDavies85/golang-film-club-app/api/internal/domain"
	"github.com/jDavies85/golang-film-club-app/api/internal/ports"
)

type SearchMoviesUC struct {
	searcher ports.MovieSearcher
	lang     string
	adults   bool
}

func NewSearchMoviesUC(searcher ports.MovieSearcher, defaultLang string, includeAdult bool) *SearchMoviesUC {
	return &SearchMoviesUC{searcher: searcher, lang: defaultLang, adults: includeAdult}
}

func (uc *SearchMoviesUC) Execute(
	ctx context.Context,
	query string,
	page int,
	language string,
	includeAdult *bool,
	year, pry *int,
) (domain.Paged[domain.Movie], error) {

	lang := uc.lang
	if language != "" {
		lang = language
	}
	adults := uc.adults
	if includeAdult != nil {
		adults = *includeAdult
	}

	return uc.searcher.SearchMovies(ctx, query, page, adults, lang, year, pry)
}
