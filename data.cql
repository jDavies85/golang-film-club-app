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

-- Insert some data
INSERT INTO filmclub.users_by_id (user_id, auth_provider, auth_subject, display_name, email, created_at)
VALUES (12345678-1234-1234-1234-123456789abc, 'entra_b2c', 'test-sub-1234', 'Test User', 'test@example.com', toTimestamp(now()));

INSERT INTO filmclub.users_by_auth (auth_provider, auth_subject, user_id)
VALUES ('entra_b2c', 'test-sub-1234', 12345678-1234-1234-1234-123456789abc);