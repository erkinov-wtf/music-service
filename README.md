# 🎵 Music Service API

Simple CRUD app, that features music service with groups and songs CRUD.

## 📋 Features

- **Complete Music Library Management**:
  - Music groups (artists/bands) CRUD operations
  - Verse-by-verse pagination
  - Preserves original formatting
  - RESTful endpoints
  - Detailed error responses
  - Filtering capabilities
  - Efficient database queries
  - Connection pooling
  - Proper error handling

## 🛠️ Tech Stack

- **Backend**: Go (Golang) with Gin web framework
- **Database**: PostgreSQL with JSONB for structured lyrics storage
- **Documentation**: Swagger/OpenAPI and Postman docs file
- **Deployment**: Docker-ready

## 📁 Project Structure

```
music-library/
├── cmd/
│   └── music-service/
│       └── main.go                 # Application entry point
├── configs/
│   ├── config.yml                  # Base configuration
│   ├── database/
│   │   ├── query.sql               # SQLc queries
│   │   ├── schema.sql              # Database schema
│   │   └── sqlc.yml                # SQLc configuration
│   ├── local.yml                   # Local environment config
│   └── release.yml                 # Production environment config
├── migrations/                     # Migration files
├── internal/
│   ├── api/
│   │   ├── handlers/               # HTTP request handlers
│   │   │   ├── group_handler.go
│   │   │   └── song_handler.go
│   │   ├── routes/                 # API routes
│   │   │   ├── path/               # Route definitions by domain
│   │   │   │   ├── group.go
│   │   │   │   └── song.go
│   │   │   ├── router.go           # Router setup
│   │   │   └── routes.go           # Route registration
│   │   └── services/               # Business logic
│   │       ├── group_service.go
│   │       └── song_service.go
│   ├── config/                     # Configuration handling
│   │   ├── config.go
│   │   └── timezone.go
│   ├── pkg/
│   │   └── utils/                  # Shared utilities
│   │       ├── constants/
│   │       │   └── timestamps.go
│   │       ├── logger/             # Structured logging
│   │       └── parser/             # Lyrics parsing utility
│   └── storage/
│       └── database/               # Database layer
│           ├── db.go               # DB connection
│           ├── models.go           # Generated models
│           ├── query.sql.go        # Generated queries
│           └── repository/         # Repository implementations
│               ├── group_repository.go
│               ├── manager.go      # Repository manager
│               └── song_repository.go
├── docs/
│   ├── docs.go                     # Swagger generated code
│   ├── music-service-postman.json  # Postman collection
│   ├── swagger.json                # OpenAPI specification
│   └── swagger.yaml                # OpenAPI YAML format
├── Dockerfile                      # Docker build instructions
├── docker-compose.yml              # Docker compose configuration
├── atlas.hcl                       # Atlas migrations configuration
```

## 🚀 Getting Started

### Installation

1. Clone the repository:

```bash
git clone https://github.com/erkinov-wtf/music-service.git
cd music-service
```

2. Install dependencies:

```bash
go mod tidy
```

3. Set up environment variables:

```bash
cp .env.docker.example .env
# Edit .env with your database credentials and other settings
```

4. Run the application:

```bash
go run cmd/music-service/main.go
```

### Docker Setup

```bash
docker compose up -d
```

## 📖 API Documentation

API documentation is available via Swagger UI at `/swagger/index.html` when the application is running.

Or use Postman docs file located at `docs/music-service-postman.json` - just import the file to Postman.

### Key Endpoints

#### Groups

- `POST /groups` - Create a new group
- `GET /groups` - List all music groups
- `GET /groups/{id}` - Get a specific group
- `PUT /groups/{id}` - Update a group
- `DELETE /groups/{id}` - Delete a group

#### Songs

- `POST /songs` - Create a new song
- `GET /songs` - List all songs with pagination
- `GET /songs/{id}` - Get a specific song
- `GET /songs/{id}/verses` - Get paginated song lyrics by verse
- `PUT /songs/{id}` - Update a song
- `DELETE /songs/{id}` - Delete a song

## 📝 Usage Examples

### Creating a Song

```bash
curl -X POST http://localhost:8080/songs \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Midnight Serenade",
    "runtime": 241,
    "lyrics": "Under the stars we dance tonight\nWhile the moon shines so bright\nIn your eyes I see the light\nOf a thousand dreams taking flight",
    "release_date": "2023-03-15T00:00:00Z",
    "link": "https://music-streaming-service.com/songs/midnight-serenade"
  }'
```

### Retrieving Paginated Lyrics Verses

```bash
curl -X GET "http://localhost:8080/songs/ee217668-3e6f-4829-946f-7bcc5cdcc595/verses?page=1&limit=5"
```

Response:
```json
{
  "song_id": "ee217668-3e6f-4829-946f-7bcc5cdcc595",
  "page": 1,
  "limit": 5,
  "pages": 10,
  "total": 48,
  "verses": [
    "Is this the real life?",
    "Is this just fantasy?",
    "Caught in a landslide",
    "No escape from reality",
    "Open your eyes"
  ]
}
```