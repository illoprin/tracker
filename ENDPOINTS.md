# Tracker API

> ℹ️ All API endpoints starts with `/api` prefix

## Endpoints

### System

| Endpoint     | Description | Requirements |
| ------------ | ----------- | ------------ |
| POST `/ping` | Ping server |              |

### User

| Endpoint           | Description           | Requirements                      |
| ------------------ | --------------------- | --------------------------------- |
| POST `/user`       | Registration          | RegisterRequest                   |
| POST `/user/login` | Log In                | LoginRequest                      |
| GET `/user/me`     | Get current user data | Authorization Token               |
| PUT `/user`        | Update current user   | UpdateRequest Authorization Token |
| DELETE `/user`     | Delete current user   | Authorization Token               |

### Artist

| Endpoint                  | Description          | Requirements                             |
| ------------------------- | -------------------- | ---------------------------------------- |
| GET `/artist/{id}`        | Get artist           |                                          |
| POST `/artist`            | Create new artist    | Authorization Token, CreateRequest       |
| GET `/artist/my`          | Get user's artists   | Authorization Token                      |
| PUT `/artist/{id}`        | Update artist        | UpdateRequest, Authorization Token       |
| PUT `/artist/{id}/avatar` | Update artist avatar | FormData, Authorization Token, Ownership |
| DELETE `/artist/{id}`     | Delete artist        | Authorization Token, Ownership           |

### Album

| Endpoint                 | Description        | Requirements                   |
| ------------------------ | ------------------ | ------------------------------ |
| POST `/track`            | Upload new track   | Authorization Token            |
| GET `/track/{id}`        | Get track metadata |                                |
| GET `/track/{id}/stream` | Stream track       |                                |
| DELETE `/track/{id}`     | Delete track       | Authorization Token, Ownership |
| PUT `/track/{id}`        | Update track       | Authorization Token, Ownership |

### Track

| Endpoint                 | Description        | Requirements                   | Status          |
| ------------------------ | ------------------ | ------------------------------ | --------------- |
| POST `/track`            | Upload new track   | Authorization Token            |                 |
| GET `/track/{id}`        | Get track metadata |                                |                 |
| GET `/track/{id}/stream` | Stream track       |                                |                 |
| DELETE `/track/{id}`     | Delete track       | Authorization Token, Ownership | Not Implemented |
| PUT `/track/{id}`        | Update track       | Authorization Token, Ownership | Not Implemented |

### Album

| Endpoint                  | Description         | Requirements                   |
| ------------------------- | ------------------- | ------------------------------ |
| POST `/album`             | Create new album    | Authorization Token            |
| GET `/album/{id}`         | Get track metadata  |                                |
| GET `/artist/{id}/albums` | Get artist's albums |                                |
| DELETE `/album/{id}`      | Delete album        | Authorization Token, Ownership |
| PUT `/album/{id}`         | Update album        | Authorization Token, Ownership |

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
  "role": enum("Admin", "Moderator", "Customer"),
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

```json
{
  "title": String,
  "status": enum('Public', 'Hidden'),
  "year": Int,
  "genres": []String
}
```
