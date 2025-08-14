package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/gocql/gocql"
	"golang.org/x/sync/errgroup"
)

var (
	ErrEmptyName = errors.New("club name is required")
)

type ClubRepository interface {
	InsertClub(ctx context.Context, clubID gocql.UUID, name string, owner gocql.UUID, createdAt time.Time) error
}

type ClubMembersRepository interface {
	InsertMember(ctx context.Context, clubID, userID, joinID gocql.UUID, role, userDisplay string) error
}

type UserClubsRepository interface {
	InsertUserClub(ctx context.Context, userID, joinID, clubID gocql.UUID, role, clubName string) error
}

// (Optional) guard to dedupe membership; no-op impl is OK for now.
type MembershipGuards interface {
	// TryInsert returns (applied, existingJoinID, err)
	TryInsert(ctx context.Context, clubID, userID, joinID gocql.UUID) (bool, gocql.UUID, error)
}

type ClubService struct {
	clubs     ClubRepository
	members   ClubMembersRepository
	userClubs UserClubsRepository
	guards    MembershipGuards // can be nil
}

func NewClubService(c ClubRepository, m ClubMembersRepository, u UserClubsRepository, g MembershipGuards) *ClubService {
	return &ClubService{clubs: c, members: m, userClubs: u, guards: g}
}

type CreateClubInput struct {
	OwnerID   gocql.UUID
	Name      string
	OwnerName string           // denormalized into members table
	Now       func() time.Time // for testability; default time.Now().UTC
}

type CreateClubOutput struct {
	ClubID    gocql.UUID
	CreatedAt time.Time
}

func (s *ClubService) Create(ctx context.Context, in CreateClubInput) (CreateClubOutput, error) {
	if in.Now == nil {
		in.Now = func() time.Time { return time.Now().UTC() }
	}
	if in.Name == "" {
		return CreateClubOutput{}, ErrEmptyName
	}
	createdAt := in.Now()
	clubID := gocql.TimeUUID()
	joinID := gocql.TimeUUID()

	// 1) Write club row
	if err := s.clubs.InsertClub(ctx, clubID, in.Name, in.OwnerID, createdAt); err != nil {
		return CreateClubOutput{}, err
	}

	// 2) Optional LWT fence to dedupe membership creation
	if s.guards != nil {
		if applied, existing, err := s.guards.TryInsert(ctx, clubID, in.OwnerID, joinID); err != nil {
			return CreateClubOutput{}, err
		} else if !applied {
			joinID = existing
		}
	}

	// 3) Dual-writes in parallel: owner as member + user-club row
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return s.members.InsertMember(ctx, clubID, in.OwnerID, joinID, "owner", in.OwnerName)
	})
	g.Go(func() error {
		return s.userClubs.InsertUserClub(ctx, in.OwnerID, joinID, clubID, "owner", in.Name)
	})
	if err := g.Wait(); err != nil {
		return CreateClubOutput{}, err
	}

	return CreateClubOutput{ClubID: clubID, CreatedAt: createdAt}, nil
}
