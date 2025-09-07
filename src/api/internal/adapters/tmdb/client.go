// internal/adapters/tmdb/client.go
package tmdbadapter

import (
	"context"
	"fmt"
	"strconv"

	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/jDavies85/golang-film-club-app/api/internal/domain"
	"github.com/jDavies85/golang-film-club-app/api/internal/ports"
)

type Client struct {
	api          *tmdb.Client
	imgBase      string // e.g., https://image.tmdb.org/t/p
	sizePoster   string // e.g., w342
	sizeBackdrop string // e.g., w780
}

func New(api *tmdb.Client, imageBase, posterSize, backdropSize string) *Client {
	return &Client{
		api:          api,
		imgBase:      imageBase,
		sizePoster:   posterSize,
		sizeBackdrop: backdropSize,
	}
}

var _ ports.MovieSearcher = (*Client)(nil)

func (c *Client) SearchMovies(
	ctx context.Context,
	query string,
	page int,
	includeAdult bool,
	language string,
	year, pry *int,
) (domain.Paged[domain.Movie], error) {
	// Note: query is now a separate argument to GetSearchMovies
	opts := map[string]string{
		"page":          strconv.Itoa(max(1, page)),
		"include_adult": strconv.FormatBool(includeAdult),
	}
	if language != "" {
		opts["language"] = language
	}
	if year != nil {
		opts["year"] = strconv.Itoa(*year)
	}
	if pry != nil {
		opts["primary_release_year"] = strconv.Itoa(*pry)
	}

	res, err := c.api.GetSearchMovies(query, opts)
	if err != nil {
		return domain.Paged[domain.Movie]{}, err
	}

	out := domain.Paged[domain.Movie]{
		Page:         res.Page,
		TotalPages:   res.TotalPages,
		TotalResults: res.TotalResults,
		Results:      make([]domain.Movie, 0, len(res.Results)),
	}
	for _, m := range res.Results {
		out.Results = append(out.Results, domain.Movie{
			ID:          int64(m.ID),
			Title:       m.Title,
			Overview:    m.Overview,
			ReleaseDate: m.ReleaseDate,
			PosterURL:   c.imageURL(m.PosterPath, c.sizePoster),
			BackdropURL: c.imageURL(m.BackdropPath, c.sizeBackdrop),
			VoteAverage: m.VoteAverage,
			VoteCount:   m.VoteCount,
		})
	}
	return out, nil
}

func (c *Client) imageURL(path string, size string) *string {
	if path == "" {
		return nil
	}
	u := fmt.Sprintf("%s/%s%s", c.imgBase, size, path)
	return &u
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
