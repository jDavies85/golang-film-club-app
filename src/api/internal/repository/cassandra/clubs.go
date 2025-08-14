package cassandra

import (
	"context"
	"time"

	"github.com/gocql/gocql"
)

type ClubRepo struct {
	insert *gocql.Query
}

func NewClubRepo(s *gocql.Session) *ClubRepo {
	return &ClubRepo{
		insert: s.Query(`INSERT INTO filmclub.film_clubs_by_id
			(club_id, name, owner_user_id, created_at) VALUES (?,?,?,?)`),
	}
}

func (r *ClubRepo) InsertClub(ctx context.Context, clubID gocql.UUID, name string, owner gocql.UUID, createdAt time.Time) error {
	return r.insert.Bind(clubID, name, owner, createdAt).WithContext(ctx).Exec()
}
