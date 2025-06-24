# Tracker API

> ℹ️ !!! All API endpoints starts with `/api/v1` prefix

## Endpoints

### System

| Endpoint    | Description |
| ----------- | ----------- |
| GET `/ping` | Ping server |

### Search

| Endpoint                                 | Description                    |
| ---------------------------------------- | ------------------------------ |
| GET `/search?query=&limitTracks=&limit=` | Search tracks, albums, artists |

### Genre

| Endpoint      | Description            |
| ------------- | ---------------------- |
| GET `/genres` | Get all allowed genres |

### Auth

| Endpoint              | Description                     | Requirements                        |
| --------------------- | ------------------------------- | ----------------------------------- |
| POST `/auth/register` | Registration                    | RegisterRequest                     |
| POST `/auth/login`    | Login                           | LoginRequest                        |
| POST `/auth/verify`   | Check status of current session | Token in JSON ({"token": "string"}) |
| POST `/auth/refresh`  | Refresh current session         | Token in JSON ({"token": "string"}) |

### User

| Endpoint          | Description                              | Requirements                 |
| ----------------- | ---------------------------------------- | ---------------------------- |
| GET `/user/me`    | Get current user                         | Authorization                |
| PATCH `/user/me`  | Update current user                      | Authorization, UpdateRequest |
| DELETE `/user/me` | Delete current user and all related data | Authorization                |

### Artist

> ❗ GET `/artist/{id}/albums` if not owning -> returns only public and approved albums

| Endpoint                               | Description                                                  | Requirements                                 |
| -------------------------------------- | ------------------------------------------------------------ | -------------------------------------------- |
| POST `/artists`                        | Create new artist                                            | Authorization, ArtistCreateRequest           |
| GET `/artists/my`                      | Get my artists                                               | Authorization                                |
| GET `/artists/liked`                   | Get liked artists                                            | Authorization                                |
| GET `/artists/{id}`                    | Get artist info                                              | Authorization                                |
| POST `/artists/{id}/like`              | Like artist                                                  | Authorization                                |
| GET `/artists/{id}/wave?limit=?&page=` | Build recommended tracks for artist (returns array TrackDTO) | Authorization                                |
| GET `/artists/{id}/stats`              | Get listening statistics                                     | Authorization                                |
| GET `/artists/{id}/albums`             | Get artist's albums                                          | Authorization                                |
| POST `/artists/{id}/album`             | Create new album for artist                                  | Authorization, Ownership, AlbumCreateRequest |
| GET `/artists/{id}/popular?limit=`     | Get popular tracks                                           | Authorization                                |
| DELETE `/artists/{id}`                 | Delete artist and related data                               | Authorization, Ownership                     |

### Genres

| Endpoint      | Description        |
| ------------- | ------------------ |
| GET `/genres` | Get allowed genres |

### Album/Track

| Endpoint                             | Description                                                  | Requirements                                 |
| ------------------------------------ | ------------------------------------------------------------ | -------------------------------------------- |
| GET `/albums/liked`                  | Get liked albums                                             | Authorization                                |
| GET `/albums/{id}`                   | Get album data                                               | Authorization                                |
| POST `/albums/{id}/like`             | Like album                                                   | Authorization                                |
| GET `/albums/{id}/listens`           | Get album listening stats                                    | Authorization                                |
| GET `/albums/{id}/wave?limit=&page=` | Build recommended tracks for album (returns array TracksDTO) | Authorization                                |
| PATCH `/albums/{id}`                 | Update album                                                 | Authorization, Ownership, AlbumUpdateRequest |
| DELETE `/albums/{id}`                | Delete album and tracks                                      | Authorization, Ownership                     |
| POST `/albums/{id}/publish`          | Send album to moderation                                     | Authorization, Ownership                     |

> ℹ️ На фронте есть две кнопки - "Сохранить изменения" и "Отправить на модерацию". Первая ничего не делает, вторая меняет параметр isPublic=true

> 💡 Модерация заключается в том, чтобы изменить isApproved=true и добавить информацию по модерации

### Track

| Endpoint                             | Description         | Requirements                                 |
| ------------------------------------ | ------------------- | -------------------------------------------- |
| POST `/albums/{id}/tracks`           | Create new track    | Authorization, Ownership, TrackCreateRequest |
| GET `/albums/{id}/wave?limit=&page=` | Get simillar tracks | Authorization                                |
| GET `/tracks/{id}`                   | Get track metadata  | Authorization                                |
| GET `/tracks/{id}/stream`            | Get track file      | Authorization                                |
| POST `/tracks/{id}/listening`        | Register Listening  | Authorization                                |

### Moderation

| Endpoint                                                      | Description           | Requirements                                     |
| ------------------------------------------------------------- | --------------------- | ------------------------------------------------ |
| GET `/albums/moderation?limit=int&artistId=uuid&query=string` | Get unapproved albums | Authorization, Moderator Role                    |
| POST `/albums/{id}/moderate`                                  | Moderate album        | Authorization, Moderator Role, ModerationRequest |

### Playlist

| Endpoint                                  | Description                                                  | Requirements                         |
| ----------------------------------------- | ------------------------------------------------------------ | ------------------------------------ |
| POST `/playlists`                         | Create new playlist                                          | Authorization, PlaylistCreateRequest |
| GET `/playlists/{id}/tracks`              | Get playlist track list                                      | Authorization                        |
| GET `/playlists/{id}/wave`                | Get recommended tracks for playlist (returns array TrackDTO) | Authorization                        |
| GET `/playlists/my`                       | Get user's playlist                                          | Authorization                        |
| PATCH `/playlists/{id}`                   | Update playlist metadata                                     | Authorization, PlaylistUpdateRequest |
| PATCH `/playlists/{id}/tracks/{trackId}`  | Push track to playlist                                       | Authorization, Ownership             |
| DELETE `/playlists/{id}/tracks/{trackId}` | Remove track from playlist                                   | Authorization, Ownership             |
| DELETE `/playlists/{id}`                  | Delete playlist                                              | Authorization, Ownership             |

## Models

### User

#### Mongo Schema

```json
{
  "id": "uuid",
  "email": "string",
  "login": "string",
  "likedArtists": ["uuid"],
  "likedAlbums": ["uuid"],
  "likedPlaylistId": "uuid",
  "passwordHash": "string",
  "createdAt": "ISODate",
  "role": "int" // 0 - regular, 1 - moderator, 2 - admin
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

#### Mongo Schema

```json
{
  "id": "uuid",
  "ownerId": "uuid",
  "name": "string",
  "avatar": "string", // file path e.g: /public/avatars/{file}
  "banner": "string", // file path e.g: /public/avatars/{file}
  "createdAt": "ISODate"
}
```

#### Create Request

```http
name: string
avatar: image
```

#### Returning Stats

##### For regular user

```json
{
  "listens_month": "int"
}
```

##### For owner

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

#### Mongo Schema

```json
{
  "id": "uuid",
  "artistId": "uuid",
  "ownerId": "uuid",
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
  "createdAt": "ISODate",
  "updatedAt": "ISODate"
}
```

#### Create Request

```http
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

#### Moderation Request

```json
{
  "status": "approved|rejected",
  "comment": "string"
}
```

#### Album DTO Responce

```json
{
  "id": "uuid",
  "artistId": "uuid",
  "ownerId": "uuid",
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
  "createdAt": "ISODate",
  "updatedAt": "ISODate",
  "count": "int",
  "duration": "int",
  "artistName": "string",
  "artistAvatar": "string"
}
```

### Playlist

#### Mongo Schema

```json
{
  "id": "uuid",
  "ownerId": "uuid",
  "name": "string",
  "description": "string",
  "cover": "string", // path to file
  "isDefault": "bool",
  "isPublic": "bool",
  "trackIds": ["string"], // array
  "createdAt": "ISODate",
  "updatedAt": "ISODate"
}
```

#### Create Request

```http
name: string
description: string
isPublic: 1|0|nil
cover: image
```

#### Update Request

```http
name?: string
description?: string
isPublic?: 1|0|nil
cover?: image
```

#### Playlist DTO Response

```json
{
  "id": "uuid",
  "name": "string",
  "description": "string",
  "cover": "string",
  "isPublic": "bool",
  "createdAt": "ISODate",
  "updatedAt": "ISODate",
  "ownerInfo": {
    "id": "uuid",
    "name": "string"
  }
}
```

### Track

#### Mongo Schema

```json
{
  "id": "uuid",
  "ownerId": "uuid", // userId
  "albumId": "uuid",
  "name": "string",
  "genres": ["string"], // array
  "duration": "int",
  "audioFile": "string", // name of audio file
  "createdAt": "ISODate"
}
```

#### Create Request

```http
name: string
genres: string # for e.g. synthpop,indie rock
audio: file # mp3, wav, flac, m4a
duration: int
```

#### Track DTO Response

```json
{
  "id": "uuid",
  "name": "string",
  "albumId": "uuid",
  "artistId": "uuid",
  "albumName": "string",
  "artistName": "string",
  "cover": "string", // file path
  "duration": "int",
  "genres": ["string"] // array
}
```

### Listening

> Listening is counted if more than 20% of the track has been listened to

#### Mongo Schema

```json
{
  "userId": "uuid",
  "trackId": "uuid",
  "duration": "int" // in seconds
}
```
