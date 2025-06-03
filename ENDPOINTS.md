# Tracker API

> ℹ️ All API endpoints starts with `/api` prefix

## Endpoints

### System

| Endpoint    | Description | Requirements |
| ----------- | ----------- | ------------ |
| GET `/ping` | Ping server |              |

### Search

| Endpoint            | Description                    | Requirements | Status          |
| ------------------- | ------------------------------ | ------------ | --------------- |
| GET `/search?query` | Search tracks, albums, artists |              | Not Implemented |

### Genre

| Endpoint      | Description            | Requirements |
| ------------- | ---------------------- | ------------ |
| GET `/genres` | Get all allowed genres |              |

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

### Track

### Album

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
