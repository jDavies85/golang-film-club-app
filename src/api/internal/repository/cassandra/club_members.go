package cassandra

import (
	"context"

	"github.com/gocql/gocql"
)

type ClubMembersRepo struct {
	insert *gocql.Query
}

func NewClubMembersRepo(s *gocql.Session) *ClubMembersRepo {
	return &ClubMembersRepo{
		insert: s.Query(`INSERT INTO filmclub.club_members_by_club
			(club_id, user_id, joined_at, role, user_display_name)
			VALUES (?,?,?,?,?)`),
	}
}

func (r *ClubMembersRepo) InsertMember(ctx context.Context, clubID, userID, joinID gocql.UUID, role, userDisplay string) error {
	return r.insert.Bind(clubID, userID, joinID, role, userDisplay).
		WithContext(ctx).Exec()
}
