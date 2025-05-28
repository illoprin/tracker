# Tracker API

> ℹ️ All API endpoints starts with `/api` prefix

## Endpoints

### System

| Endpoint    | Description | Requirements |
| ----------- | ----------- | ------------ |
| GET `/ping` | Ping server |              |

### Search

| Endpoint                     | Description                    | Requirements | Status          |
| ---------------------------- | ------------------------------ | ------------ | --------------- |
| GET `/resource/search?query` | Search tracks, albums, artists |              | Not Implemented |

### Genre

| Endpoint      | Description            | Requirements |
| ------------- | ---------------------- | ------------ |
| GET `/genres` | Get all allowed genres |              |

### User

| Endpoint                 | Description           | Requirements                        |
| ------------------------ | --------------------- | ----------------------------------- |
| POST `/user`             | Registration          | RegisterRequest                     |
| POST `/user/login`       | Log In                | LoginRequest                        |
| GET `/user/me`           | Get current user data | Authorization Token                 |
| PUT `/user`              | Update current user   | UpdateRequest Authorization Token   |
| DELETE `/user`           | Delete current user   | Authorization Token                 |
| GET `/user/search?query` | Search users          | Authorization Token, Moderator role |

### Artist

| Endpoint                  | Description          | Requirements                             |
| ------------------------- | -------------------- | ---------------------------------------- |
| GET `/artist/{id}`        | Get artist           |                                          |
| GET `/artist/{id}/albums` | Get artist's albums  |                                          |
| POST `/artist`            | Create new artist    | Authorization Token, CreateRequest       |
| GET `/artist/my`          | Get user's artists   | Authorization Token                      |
| PUT `/artist/{id}`        | Update artist        | UpdateRequest, Authorization Token       |
| PUT `/artist/{id}/avatar` | Update artist avatar | FormData, Authorization Token, Ownership |
| DELETE `/artist/{id}`     | Delete artist        | Authorization Token, Ownership           |

### Track

| Endpoint                 | Description        | Requirements                   | Status          |
| ------------------------ | ------------------ | ------------------------------ | --------------- |
| POST `/track`            | Upload new track   | Authorization Token            |                 |
| GET `/track/{id}`        | Get track metadata |                                |                 |
| GET `/track/{id}/stream` | Stream track       | HTTP-Range request             |                 |
| DELETE `/track/{id}`     | Delete track       | Authorization Token, Ownership |                 |
| PUT `/track/{id}`        | Update track       | Authorization Token, Ownership | Not Implemented |

### Album

| Endpoint                 | Description                 | Requirements                   |
| ------------------------ | --------------------------- | ------------------------------ |
| POST `/album`            | Create new album            | Authorization Token            |
| GET `/album/{id}`        | Get album metadata          |                                |
| GET `/album/{id}/tracks` | Get album's tracks metadata |                                |
| DELETE `/album/{id}`     | Delete album                | Authorization Token, Ownership |
| PUT `/album/{id}`        | Update album                | Authorization Token, Ownership |

### Content Moderation

| Endpoint                     | Description              | Requirements                                            |
| ---------------------------- | ------------------------ | ------------------------------------------------------- |
| PUT `/album/{id}/moderation` | Moderate album           | Authorization Token, Moderator role, Moderation request |
| GET `/album/on-moderation`   | Get albums on moderation | Authorization Token, Moderator role                     |

### Playlist

| Endpoint                                 | Description                    | Requirements                                            |
| ---------------------------------------- | ------------------------------ | ------------------------------------------------------- |
| POST `/playlist/`                        | Create new playlist            | Authorization Token                                     |
| GET `/playlist/{id}/tracks`              | Get playlist's tracks metadata | Authorization Token, Ownership if resource isn't public |
| PUT `/playlist/{id}/tracks/{trackID}`    | Push track to playlist         | Authorization Token, Ownership                          |
| DELETE `/playlist/{id}/tracks/{trackID}` | Remove track from playlist     | Authorization Token, Ownership                          |
| PUT `/playlist/{id}`                     | Update playlist metadata       | Authorization Token, Ownership                          |
| DELETE `/playlist/{id}`                  | Delete playlist metadata       | Authorization Token, Ownership                          |

## Models

### User

#### Schema

```json
{
  "id": StringUUID,
  "login": String,
  "email": String,
  "passwordHash": String,
  "myChoicePlaylist": StringUUID,
  "createdAt": ISO8601Date,
  "role": enum("Admin", "Moderator", "Customer"),
}
```

#### Token

Token payload contains json

```json
{
  "id": String,
  "email": String,
  "role": Int,
}
```

#### Register

```json
{
  "login": String,
  "email": String,
  "password": String,
}
```

#### Login

```json
{
  "login": String,
  "password": String,
}
```

#### Update

```json
{
  "login"?: String,
  "password"?: String,
  "email"?: String,
  "role"?: String, // if 'Allow-Access' header set'
}
```

### Artist

#### Schema

```json
{
  "id": StringUUID,
  "name": String,
  "userID": StringUUID,
  "avatarPath": String,
  "createdAt": ISO8601Date,
}
```

#### Create

```json
{
  "name": String,
}
```

#### Update

```json
{
  "name"?: String,
}
```

#### Update avatar Form Data

```http
avatar: file
```

### Track

#### Schema

```json
{
  "id": StringUUID,
  "title": String,
  "duration": Int,
  "genres": []String,
  "audioFile": String, // file name
  "albumID": StringUUID,
  "createdAt": ISO8601Date
}
```

#### Create Form Data

```http
title: string
genre: []string
albumId: stringUUID
audio: audio/wav,audio/m4a,audio/mp3
```

### Album

#### Schema

```json
{
  "id": StringUUID,
  "title": String,
  "artistID": StringUUID,
  "year": Int,
  "coverPath": String, // path to cover file
  "genres": []String,
  "status": enum('Public', 'Hidden', 'OnModeration'),
  "createdAt": ISO8601Date
}
```

#### Create request

```json
{
  "title": String,
  "artistID": StringUUID,
  "year": Int,
  "genres": []String
}
```

#### Update request

> ℹ️ after any changes the album goes to moderation (status = 'OnModeration')

```json
{
  "title?": String,
  "isHidden?": Bool,
  "year?": Int,
  "genres?": []String
}
```

#### Moderation request

```json
{
  "status": enum('Moderated', 'Denied')
  "reason": String
}
```

### Playlist

> ℹ️ default playlist "My Choice" is marked as default (isDefault = true)

#### Schema

```json
{
  "name": String,
  "userID": StringUUID,
  "isDefault": Bool,
  "isPublic": Bool,
  "trackIDs": []StringUUID,
  "updatedAt": ISO8601Date
}
```

#### Update request

> ℹ️ each user can update non-default playlists only

```json
{
  "name?": String,
  "isPublic?": Bool,
}
```
