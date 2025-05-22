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

| Endpoint              | Description        | Requirements                       |
| --------------------- | ------------------ | ---------------------------------- |
| GET `/artist/{id}`    | Get artist         |                                    |
| POST `/artist`        | Create new artist  | Authorization Token, CreateRequest |
| GET `/artist/my`      | Get user's artists | Authorization Token                |
| PUT `/artist/{id}`    | Update artist      | UpdateRequest Authorization Token  |
| DELETE `/artist/{id}` | Delete artist      | Authorization Token                |

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
  "avatarPath"?: String,
}
```
