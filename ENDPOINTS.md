# Tracker API

> ℹ️ All API endpoints starts with `/api` prefix

## Endpoints

### System

| Endpoint    | Description |
| ----------- | ----------- |
| GET `/ping` | Ping server |

### Search

| Endpoint            | Description                    | Status          |
| ------------------- | ------------------------------ | --------------- |
| GET `/search?query` | Search tracks, albums, artists | Not Implemented |

### Genre

| Endpoint      | Description            |
| ------------- | ---------------------- |
| GET `/genres` | Get all allowed genres |

### Auth

| Endpoint              | Description                     | Requirements    |
| --------------------- | ------------------------------- | --------------- |
| POST `/auth/register` | Registration                    | RegisterRequest |
| POST `/auth/login`    | Login                           | LoginRequest    |
| POST `/auth/verify`   | Check status of current session | Token in JSON   |
| POST `/auth/refresh`  | Refresh current session         | Token in JSON   |

### User

| Endpoint          | Description                              | Requirements                 |
| ----------------- | ---------------------------------------- | ---------------------------- |
| GET `/user/me`    | Get current user                         | Authorization                |
| PATCH `/user/me`  | Update current user                      | Authorization, UpdateRequest |
| DELETE `/user/me` | Delete current user and all related data | Authorization                |

### Artist

> ❗ GET `/artist/{id}/albums` if not owning -> returns only public albums

| Endpoint                           | Description                    | Requirements                       |
| ---------------------------------- | ------------------------------ | ---------------------------------- |
| POST `/artists`                    | Create new artist              | Authorization, ArtistCreateRequest |
| GET `/artists/my`                  | Get my artists                 | Authorization                      |
| GET `/artists/{id}`                | Get artist info                | Authorization                      |
| GET `/artists/{id}/stats`          | Get listening statistics       | Authorization                      |
| GET `/artists/{id}/albums`         | Get artist's albums            | Authorization                      |
| GET `/artists/{id}/popular?limit=` | Get popular tracks             | Authorization                      |
| DELETE `/artists/{id}`             | Delete artist and related data | Authorization, Ownership           |

### Genres

| Endpoint      | Description        | Requirements |
| ------------- | ------------------ | ------------ |
| GET `/genres` | Get allowed genres |              |

### Album/Track

| Endpoint                       | Description               | Requirements                                 |
| ------------------------------ | ------------------------- | -------------------------------------------- |
| POST `/api/albums`             | Create new album          | Authorization, AlbumCreateRequest            |
| POST `/api/albums/{id}/tracks` | Create new track          | Authorization, TrackCreateRequest            |
| GET `/api/albums/{id}`         | Get album data            | Authorization                                |
| GET `/api/albums/{id}/listens` | Get album listening stats | Authorization                                |
| GET `/api/albums/{id}/tracks`  | Get album tracks          | Authorization                                |
| PATCH `/api/albums/{id}`       | Update album              | Authorization, Ownership, AlbumUpdateRequest |
| DELETE `/api/albums/{id}`      | Delete album and tracks   | Authorization, Ownership                     |

#### Listening

### Moderation

### Playlist

## Models

### User

```json
{
  "id": "id",
  "email": "string",
  "login": "string",
  "likedArtists": ["id"],
  "likedAlbums": ["id"],
  "likedPlaylistId": "id",
  "passwordHash": "string",
  "createdAt": "ISODate",
  "role": "int" // 0 - обычный, 1 - модератор, 2 - админ
}
```

#### RegisterRequest

```json
{
  "email": "string",
  "login": "string",
  "password": "string"
}
```

#### LoginRequest

```json
{
  "login": "string",
  "password": "string"
}
```

#### UpdateRequest

```json
{
  "login": "string",
  "password": "string",
  "email": "string"
}
```

### Artist

```json
{
  "id": "id",
  "ownerId": "id",
  "name": "string",
  "avatar": "string", // file path e.g: /public/avatars/{file}
  "createdAt": "ISODate"
}
```

#### Create Request

```http
name: string
avatar: image
```

#### Stats

##### Regular user

```json
{
  "listens_month": "int"
}
```

##### Owner

```json
{
  "listens_month": "int",
  "albums": [
    {
      "id": "id",
      "listens_month": "int"
    }
  ]
}
```

### Album

```json
{
  "id": "id",
  "artistId": "id",
  "ownerId": "id",
  "name": "string",
  "year": "int",
  "cover": "string",
  "type": "single|album", 
  "isPublic": "bool",
  "isApproved": "bool",
  "moderation": {
    "status": "pending|approved|rejected",
    "comment": "string"
  },
  "createdAt": "ISODate"
}
```

#### Create Request

```http
artistId: id
name: string
year: int
type: single|album
cover: image
```

#### Update Request

> ℹ️ isPublic hides album in global search and artist's profile, prevents moderation

```http
name: string
year: int
type: single|album
isPublic: bool
```