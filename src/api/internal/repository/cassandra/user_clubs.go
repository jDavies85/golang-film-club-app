package cassandra

import (
	"context"

	"github.com/gocql/gocql"
)

type UserClubsRepo struct {
	insert *gocql.Query
}

func NewUserClubsRepo(s *gocql.Session) *UserClubsRepo {
	return &UserClubsRepo{
		insert: s.Query(`INSERT INTO filmclub.user_clubs_by_user
			(user_id, joined_at, club_id, role, club_name)
			VALUES (?,?,?,?,?)`),
	}
}

func (r *UserClubsRepo) InsertUserClub(ctx context.Context, userID, joinID, clubID gocql.UUID, role, clubName string) error {
	return r.insert.Bind(userID, joinID, clubID, role, clubName).
		WithContext(ctx).Exec()
}
