package domain

type Movie struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"releaseDate"`
	PosterURL   *string `json:"posterUrl,omitempty"`
	BackdropURL *string `json:"backdropUrl,omitempty"`
	VoteAverage float32 `json:"voteAverage"`
	VoteCount   int64   `json:"voteCount"`
}

type Paged[T any] struct {
	Page         int64 `json:"page"`
	TotalPages   int64 `json:"totalPages"`
	TotalResults int64 `json:"totalResults"`
	Results      []T
}
