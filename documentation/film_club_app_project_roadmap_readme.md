# Film Club App – Project Roadmap

A pragmatic, Go-idiomatic roadmap you can track directly in GitHub. It’s organized into phases sized for evenings/weekends, with checkboxes, reasoning, code snippets, and CQL.

> **Stack**: Go + Gin (API), Cassandra (Cosmos DB Cassandra API later), gocql (driver), React (Vite), Azure Static Web Apps (FE), Azure App Service (API), Azure Blob Storage (images), Azure Entra B2C (auth), TMDB for search.

---

## ✅ Getting Started

- [ ] Create repo and add this README
- [ ] Initialize Go module: `go mod init github.com/you/filmclub`
- [ ] Add `.gitignore` (Go, Node, env files)
- [ ] Create `.env.example` and (local) `.env` (git-ignored)

```env
# .env.example
TMDB_API_KEY=
APP_HTTP_PORT=8080
APP_CASSANDRA_HOSTS=127.0.0.1
APP_CASSANDRA_KEYSPACE=filmclub
APP_CASSANDRA_CONSISTENCY=LOCAL_QUORUM
```

---

## Phase 1 — Skeleton API & Local Cassandra

**Goal:** Run a health endpoint; talk to local Cassandra.

- [ ] Project layout
```
filmclub/
├─ apps/
│  ├─ api/                    # Golang API (its own go.mod)
│  │  ├─ cmd/api/main.go
│  │  ├─ internal/
│  │  │  ├─ http/             # router, middleware, handlers
│  │  │  ├─ usecase/          # services / CQRS-light
│  │  │  ├─ repository/       
│  │  │  │  └─ cassandra/     # session + repos
│  │  │  ├─ external/
│  │  │  │  └─ tmdb/          # wrapper over github.com/cyruzin/golang-tmdb
│  │  │  ├─ config/           
│  │  │  └─ observability/    # logging, tracing, metrics (optional)
│  │  ├─ migrations/          # *.cql (with filmclub.<table> keyspace-qualified)
│  │  ├─ scripts/             # local helper scripts (seed, smoke, etc.)
│  │  ├─ .env.example
│  │  └─ go.mod
│  └─ web/                    # React app (Vite or CRA)
│     ├─ src/
│     ├─ public/
│     ├─ .env.example         # VITE_* vars only
│     ├─ package.json
│     └─ vite.config.ts
├─ infra/
│  ├─ docker-compose.yml      # local dev: cassandra (+ optional admin UIs)
│  ├─ cassandra/
│  │  └─ data.cql             # seed data; safe to re-run (idempotent inserts)
│  ├─ azure/
│  │  ├─ appservice/          # Bicep/Terraform (optional)
│  │  └─ staticwebapps/       # swa config (routes.json), auth, etc.
│  └─ k8s/                    # if you ever containerize (optional)
├─ docs/
│  └─ README.md               # your roadmap/checklist lives here
├─ .github/
│  └─ workflows/
│     ├─ api-ci.yml           # go build, test, lint; deploy on main (optional)
│     └─ web-ci.yml           # web build, lint; deploy SWA (optional)
├─ Makefile                   # convenience entry points for the whole repo
├─ .env.example               # root-level shared defaults (non-secrets)
├─ .gitignore
└─ .editorconfig
```

- [ ] Add dependencies: `gin-gonic/gin`, `gocql/gocql`, `joho/godotenv`, `golang.org/x/sync/errgroup`, `golang.org/x/sync/singleflight`
- [ ] Health endpoint

```go
// cmd/api/main.go (excerpt)
_ = godotenv.Load() // local only; prod uses env vars
r := gin.Default()
r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
r.Run(":" + os.Getenv("APP_HTTP_PORT"))
```

- [ ] Local Cassandra via Docker

```bash
docker run --name cassandra -p 9042:9042 -d cassandra:4.1
docker exec -it cassandra cqlsh
```

- [ ] Keyspace (dev):
```sql
CREATE KEYSPACE IF NOT EXISTS filmclub
WITH REPLICATION = { 'class': 'SimpleStrategy', 'replication_factor': '1' };
```

**Reasoning:** Keep the skeleton tiny to minimize startup friction. Local Cassandra first; switch to Cosmos later.

---

## Phase 2 — Users (Auth Mapping) & Config

**Goal:** Register/log first user using Entra B2C `sub` claim → create internal `user_id`.

- [ ] Config loader
```go
// internal/config/config.go
package config
import "github.com/kelseyhightower/envconfig"

type Config struct {
  HTTP struct{ Port string `default:"8080"` }
  Cassandra struct{
    Hosts        []string `envconfig:"HOSTS" required:"true"`
    Keyspace     string   `required:"true"`
    Consistency  string   `default:"LOCAL_QUORUM"`
  }
  TMDB struct{ APIKey string `envconfig:"API_KEY"` }
}

func Load() (Config, error) { var c Config; return c, envconfig.Process("APP", &c) }
```

- [ ] CQL tables
```sql
-- users
CREATE TABLE IF NOT EXISTS users_by_id (
  user_id uuid PRIMARY KEY,
  auth_provider text,
  auth_subject text,
  display_name text,
  email text,
  created_at timestamp
);

CREATE TABLE IF NOT EXISTS users_by_auth (
  auth_provider text,
  auth_subject text,
  user_id uuid,
  PRIMARY KEY ((auth_provider), auth_subject)
);
```

- [ ] Sample test user (insert in `data.cql`)
```sql
INSERT INTO users_by_id (user_id, auth_provider, auth_subject, display_name, email, created_at)
VALUES (uuid(), 'entra_b2c', 'test-sub-1234', 'Test User', 'test@example.com', toTimestamp(now()));

INSERT INTO users_by_auth (auth_provider, auth_subject, user_id)
VALUES ('entra_b2c', 'test-sub-1234', <same-uuid-as-above>);
```
*(Tip: for deterministic UUIDs, replace `uuid()` with a fixed UUID string you generate once, e.g. `12345678-1234-1234-1234-123456789abc` so inserts are idempotent when rerunning.)*

**Reasoning:** Separate internal `user_id` from external identity to future-proof provider changes. Having a test user bootstraps local dev.

---

## Phase 3 — Film Clubs & Membership (Dual Writes + errgroup)

**Goal:** Create club; owner becomes first member. Dual writes with retries & idempotency.

- [ ] CQL tables
```sql
CREATE TABLE IF NOT EXISTS film_clubs_by_id (
  club_id uuid PRIMARY KEY,
  name text,
  owner_user_id uuid,
  created_at timestamp
);

CREATE TABLE IF NOT EXISTS user_clubs_by_user (
  user_id uuid,
  joined_at timeuuid,
  club_id uuid,
  role text,
  club_name text,
  PRIMARY KEY ((user_id), joined_at, club_id)
) WITH CLUSTERING ORDER BY (joined_at DESC);

CREATE TABLE IF NOT EXISTS club_members_by_club (
  club_id uuid,
  user_id uuid,
  joined_at timeuuid,
  role text,
  user_display_name text,
  PRIMARY KEY ((club_id), user_id)
);

-- optional LWT guard to dedupe membership
CREATE TABLE IF NOT EXISTS membership_guards (
  club_id uuid,
  user_id uuid,
  join_id timeuuid,
  PRIMARY KEY ((club_id), user_id)
);
```

- [ ] Usecase with `errgroup` and idempotent `join_id`
```go
func (s *ClubService) AddMember(ctx context.Context, clubID, userID gocql.UUID, role, clubName, userName string) error {
  joinID := gocql.TimeUUID() // generate once per op
  // optional LWT fence
  applied, existing, err := s.guards.TryInsert(ctx, clubID, userID, joinID)
  if err != nil { return err }
  if !applied { joinID = existing }

  writeOnce := func(ctx context.Context) error {
    g, ctx := errgroup.WithContext(ctx)
    g.Go(func() error { return s.members.Insert(ctx, clubID, userID, joinID, role, userName) })
    g.Go(func() error { return s.userClubs.Insert(ctx, userID, joinID, clubID, role, clubName) })
    return g.Wait()
  }

  for i := 0; i < 3; i++ {
    cctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    err = writeOnce(cctx)
    cancel()
    if err == nil { return nil }
    time.Sleep(time.Duration(200*(i+1)) * time.Millisecond)
  }
  // enqueue background repair (outbox)
  return err
}
```

- [ ] Endpoints
  - [ ] `POST /v1/clubs` → returns `club_id`
  - [ ] `GET /v1/me/clubs`
  - [ ] `GET /v1/clubs/{club_id}/members`

**Reasoning:** Cassandra has no cross-partition transactions; we fix-forward using idempotent calls and optional LWT guard.

---

## Phase 4 — Invites (Codes + Email UX)

**Goal:** Owner generates invite codes; invitees accept to join.

- [ ] CQL tables
```sql
CREATE TABLE IF NOT EXISTS invites_by_code (
  code text PRIMARY KEY,
  club_id uuid,
  issuer_user_id uuid,
  created_at timestamp,
  expires_at timestamp
);

CREATE TABLE IF NOT EXISTS invites_by_email (
  email text,
  created_at timeuuid,
  club_id uuid,
  code text,
  club_name text,
  issuer_user_id uuid,
  PRIMARY KEY ((email), created_at)
) WITH CLUSTERING ORDER BY (created_at DESC);
```

- [ ] Endpoints
  - [ ] `POST /v1/clubs/{club_id}/invites { emails?: [] }` → returns `{ code }` and/or `codesByEmail`
  - [ ] `POST /v1/invites/{code}/accept` → calls `AddMember`

**Reasoning:** Two views give O(1) lookup for “resolve code” and “my pending invites”.

---

## Phase 5 — TMDB Search & Movie Selection

**Goal:** Search TMDB; store a club’s current pick with deadline.

- [ ] TMDB adapter and service
```go
// internal/usecase/movie_service.go
type Movie struct{ ID int `json:"id"`; Title string `json:"title"`; Year int `json:"year"` }

type MovieSearchService interface{ Search(ctx context.Context, q string) ([]Movie, error) }

type MovieService struct{ searcher MovieSearchService }
func NewMovieService(s MovieSearchService) *MovieService { return &MovieService{searcher: s} }
func (svc *MovieService) SearchMovies(ctx context.Context, q string) ([]Movie, error) {
  return svc.searcher.Search(ctx, q)
}
```

```go
// internal/external/tmdb/client.go
api, _ := tmdb.Init(os.Getenv("TMDB_API_KEY"))
// implement Search(ctx, q) to map results → usecase.Movie
```

- [ ] CQL for selections (history + current)
```sql
CREATE TABLE IF NOT EXISTS club_movies_by_club (
  club_id uuid,
  cycle_id timeuuid,
  movie_id int,
  title text,
  chosen_by_user_id uuid,
  chosen_at timestamp,
  poster_url text,
  deadline timestamp,
  status text,
  PRIMARY KEY ((club_id), cycle_id)
) WITH CLUSTERING ORDER BY (cycle_id DESC);
```

- [ ] Endpoints
  - [ ] `GET /v1/movies/search?q=…`
  - [ ] `POST /v1/clubs/{club_id}/selections` (sets `cycle_id = now()`)
  - [ ] `GET /v1/clubs/{club_id}/selections/current`

**Reasoning:** Time-series pattern with latest-first; denormalize poster/title for snappy UI.

---

## Phase 6 — React Seed UI

**Goal:** Minimal UI to exercise flows.

- [ ] Vite React app scaffold
- [ ] Pages
  - [ ] Register (implicit on first auth)
  - [ ] Create first club
  - [ ] Invite members
  - [ ] Search & select movie; set deadline
- [ ] `.env` proxy to API; fetch with bearer token later

**Reasoning:** Vertical slices keep momentum: working feature end-to-end beats breadth.

---

## Phase 7 — Auth (Azure Entra B2C)

**Goal:** Protect mutating endpoints.

- [ ] Register SPA + API apps in B2C
- [ ] React: add MSAL.js for login and token acquisition
- [ ] API: JWT middleware (validate issuer, audience, keys)
- [ ] Middleware: map `sub` → `user_id` using `users_by_auth`; auto-create if missing

**Reasoning:** Keep GETs public initially; require auth for POST/PUT/DELETE.

---

## Phase 8 — Media Uploads (Azure Blob Storage)

**Goal:** Users attach small images to selections/reviews.

- [ ] Create storage account + container
- [ ] API: generate SAS upload URL for authenticated user
- [ ] React: upload directly to Blob using SAS; store blob URL in Cassandra

**Reasoning:** Direct-to-blob avoids API buffering; store only the URL in DB.

---

## Phase 9 — Observability & Quality

- [ ] Structured logging (request id, user id, route, latency)
- [ ] Health/ready endpoints for Azure
- [ ] Linting: `golangci-lint`
- [ ] Unit tests (usecase with fakes/mocks)
- [ ] Integration tests with Docker Cassandra in `test/`

**Reasoning:** Early visibility prevents time sinks later.

---

## Phase 10 — Azure Deployment (free/low cost)

- [ ] Deploy React to Azure Static Web Apps (GitHub integration)
- [ ] Deploy API to Azure App Service (Linux) or container
- [ ] App settings: set env vars (TMDB key, Cassandra hosts/keyspace, B2C settings)
- [ ] (Later) Move DB to Cosmos DB Cassandra API (free tier) and update env vars

**Reasoning:** Keep infra simple; promote from local → cloud with same env names.

---

## Migrations Strategy (Cassandra)

- [ ] Keep idempotent `.cql` under `migrations/`
- [ ] Add `cmd/migrate` or use `gocqlx/migrate` / `golang-migrate`

```sql
-- migrations/001_keyspace.cql
CREATE KEYSPACE IF NOT EXISTS filmclub WITH REPLICATION = { 'class': 'SimpleStrategy', 'replication_factor': '1' };

-- migrations/010_users.cql
CREATE TABLE IF NOT EXISTS users_by_id (...);
CREATE TABLE IF NOT EXISTS users_by_auth (...);

-- migrations/020_clubs.cql
CREATE TABLE IF NOT EXISTS film_clubs_by_id (...);
CREATE TABLE IF NOT EXISTS user_clubs_by_user (...);
CREATE TABLE IF NOT EXISTS club_members_by_club (...);
CREATE TABLE IF NOT EXISTS membership_guards (...);

-- migrations/030_invites.cql
CREATE TABLE IF NOT EXISTS invites_by_code (...);
CREATE TABLE IF NOT EXISTS invites_by_email (...);

-- migrations/040_selections.cql
CREATE TABLE IF NOT EXISTS club_movies_by_club (...);
```

**Reasoning:** Idempotent CQL enables replays and easy CI automation.

---

## Concurrency & Go-isms Cheatsheet

- **errgroup** for parallel writes (dual views)
- **singleflight** to dedupe hot searches
- **context deadlines** around I/O
- **background ticker** for deadline reminders
- **SSE/WebSocket hub** for live updates

```go
if err := g.Wait(); err != nil { /* first error from goroutines */ }
```

**Reasoning:** Embrace explicit concurrency with safe cancellation and retries.

---

## API Sketch (first endpoints)

- [ ] `GET    /health`
- [ ] `POST   /v1/users/me/register`
- [ ] `POST   /v1/clubs`
- [ ] `GET    /v1/me/clubs`
- [ ] `GET    /v1/clubs/{club_id}/members`
- [ ] `POST   /v1/clubs/{club_id}/invites`
- [ ] `POST   /v1/invites/{code}/accept`
- [ ] `GET    /v1/movies/search?q=`
- [ ] `POST   /v1/clubs/{club_id}/selections`
- [ ] `GET    /v1/clubs/{club_id}/selections/current`

---

## Backlog / Nice-to-haves

- [ ] Pagination cursors (use last clustering key)
- [ ] Reviews/comments per selection (time-series table)
- [ ] Notifications (email/Teams) before deadlines
- [ ] Role-based permissions beyond owner/member
- [ ] Rate limiting for TMDB and API
- [ ] Cache TMDB lookups in Cassandra or memory (short TTL)

---

## Running Locally

```bash
# API
APP_HTTP_PORT=8080 APP_CASSANDRA_HOSTS=127.0.0.1 APP_CASSANDRA_KEYSPACE=filmclub \
APP_TMDB_API_KEY=xxxx go run ./cmd/api

# React (separate repo or /web)
npm create vite@latest web -- --template react-ts
cd web && npm install && npm run dev
```

---

## Notes on Cosmos DB (Cassandra API)

- Start locally with OSS Cassandra; move to Cosmos later.
- Keep partition keys high-cardinality (`user_id`, `club_id`), avoid hot partitions.
- Use idempotent writes; watch RU/s on scans; prefer key lookups.

---

## License & Attribution

- TMDB terms require attribution if you display their data. Include their logo & link where appropriate.

---

Happy building! Check off items as you go and adjust phases to fit your time. This roadmap aims for steady momentum with clean, Go-idiomatic slices.

