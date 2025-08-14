-- Create a keyspace
CREATE KEYSPACE IF NOT EXISTS filmclub
WITH REPLICATION = { 'class': 'SimpleStrategy', 'replication_factor': '1' };

-- Create a table
CREATE TABLE IF NOT EXISTS filmclub.users_by_id (
  user_id uuid PRIMARY KEY,
  auth_provider text,
  auth_subject text,
  display_name text,
  email text,
  created_at timestamp
);

CREATE TABLE IF NOT EXISTS filmclub.users_by_auth (
  auth_provider text,
  auth_subject text,
  user_id uuid,
  PRIMARY KEY ((auth_provider), auth_subject)
);

CREATE TABLE IF NOT EXISTS filmclub.film_clubs_by_id (
  club_id uuid PRIMARY KEY,
  name text,
  owner_user_id uuid,
  created_at timestamp
);

CREATE TABLE IF NOT EXISTS filmclub.user_clubs_by_user (
  user_id uuid,
  joined_at timeuuid,
  club_id uuid,
  role text,
  club_name text,
  PRIMARY KEY ((user_id), joined_at, club_id)
) WITH CLUSTERING ORDER BY (joined_at DESC);

CREATE TABLE IF NOT EXISTS filmclub.club_members_by_club (
  club_id uuid,
  user_id uuid,
  joined_at timeuuid,
  role text,
  user_display_name text,
  PRIMARY KEY ((club_id), user_id)
);

-- optional LWT guard to dedupe membership
CREATE TABLE IF NOT EXISTS filmclub.membership_guards (
  club_id uuid,
  user_id uuid,
  join_id timeuuid,
  PRIMARY KEY ((club_id), user_id)
);

-- Insert some data
INSERT INTO filmclub.users_by_id (user_id, auth_provider, auth_subject, display_name, email, created_at)
VALUES (12345678-1234-1234-1234-123456789abc, 'entra_b2c', 'test-sub-1234', 'Test User', 'test@example.com', toTimestamp(now()));

INSERT INTO filmclub.users_by_auth (auth_provider, auth_subject, user_id)
VALUES ('entra_b2c', 'test-sub-1234', 12345678-1234-1234-1234-123456789abc);