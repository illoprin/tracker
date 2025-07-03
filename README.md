[Russian Version](README.ru.md) | [English Version](README.md)

# 🎸 Tracker - Music Streaming Service

## 📃 Table of Contents

- [Project Idea](#idea)
- [Key Features](#features)
- [Result](#result)
- [Local Setup](#run)

<a name="idea"></a> 

## 💡 Idea

**Tracker** is a chronicle of personal music preferences. An app for so-called "musical introverts," where anyone can expand the existing music library.

<a name="features"></a> 

## ✨ Key Features

- **Recommendation System**  
  
  All tracks on the homepage, artist pages, and albums are selected and sorted to match the musical preferences of each individual user.

- **No Boundaries Between Listener and Artist**  
  
  Anyone can upload their songs by registering an artist profile. Users gain access to detailed statistics on plays, audience, and the popularity of their tracks.

<a name="result"></a> 

## ❗ Result

As of now, the API functionality has been fully implemented.

You can explore the endpoints either in [SwaggerUI](localhost:8000/api/v1/docs) or in the [ENDPOINTS.md](https://github.com/illoprin/tracker/blob/master/ENDPOINTS.md) file.

<a name="run"></a> 

## ▶️ Local Setup

You can run the web app in development mode (with hot-reload) using Docker:

```bash
git clone https://github.com/illoprin/tracker.git  
cd tracker  
docker compose up -d --build  
```

The API will be available at **localhost:8000**