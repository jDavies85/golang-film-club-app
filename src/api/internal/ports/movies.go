package ports

import (
	"context"

	"github.com/jDavies85/golang-film-club-app/api/internal/domain"
)

type MovieSearcher interface {
	SearchMovies(
		ctx context.Context,
		query string,
		page int,
		includeAdult bool,
		language string,
		year, primaryReleaseYear *int,
	) (domain.Paged[domain.Movie], error)
}
