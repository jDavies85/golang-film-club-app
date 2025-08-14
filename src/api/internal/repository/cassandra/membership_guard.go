package cassandra

import (
	"context"

	"github.com/gocql/gocql"
)

type MembershipGuardRepo struct {
	insertLWT *gocql.Query
	selectOne *gocql.Query
}

func NewMembershipGuardRepo(s *gocql.Session) *MembershipGuardRepo {
	return &MembershipGuardRepo{
		insertLWT: s.Query(`INSERT INTO filmclub.membership_guards (club_id, user_id, join_id)
			VALUES (?,?,?) IF NOT EXISTS`),
		selectOne: s.Query(`SELECT join_id FROM filmclub.membership_guards WHERE club_id=? AND user_id=?`),
	}
}

func (r *MembershipGuardRepo) TryInsert(ctx context.Context, clubID, userID, joinID gocql.UUID) (bool, gocql.UUID, error) {
	m := map[string]interface{}{}
	applied, err := r.insertLWT.Bind(clubID, userID, joinID).
		WithContext(ctx).MapScanCAS(m)
	if err != nil {
		return false, gocql.UUID{}, err
	}
	if applied {
		return true, joinID, nil
	}
	// fetch existing join_id
	var existing gocql.UUID
	if err := r.selectOne.Bind(clubID, userID).WithContext(ctx).Scan(&existing); err != nil {
		return false, gocql.UUID{}, err
	}
	return false, existing, nil
}
