package domain

type Movie struct {
	ID          int64
	Title       string
	Overview    string
	ReleaseDate string
	PosterURL   *string
	BackdropURL *string
	VoteAverage float32
	VoteCount   int64
}

type Paged[T any] struct {
	Page         int64
	TotalPages   int64
	TotalResults int64
	Results      []T
}
