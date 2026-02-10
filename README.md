[Russian Version](README.ru.md) | [English Version](README.md)

# 🎸 Tracker - Music Streaming Service

## 📃 Table of Contents

- [Screenshots](#screenshots)
- [Project Idea](#idea)
- [Key Features](#features)
- [Result](#result)
- [Local Setup](#run)

<a name="screenshots"></a> 

## 🖼️ Screenshots

**Home page**

![Image](https://github.com/user-attachments/assets/acf37b82-592b-4848-becf-9f470288e478)

**Registration form**

![Image](https://github.com/user-attachments/assets/f12d4ade-f6de-4833-a017-51a6e2e8e049)

**Artist page**

![Image](https://github.com/user-attachments/assets/5c62f55c-787f-4a72-a071-8d4656b18a6c)

**Album page**

![Image](https://github.com/user-attachments/assets/ae96145a-32bf-49a0-a640-f94d3654083e)

**New artist form**

![Image](https://github.com/user-attachments/assets/9dad0157-78c0-4d27-8c99-a720f1e07af8)

**New album form**

![Image](https://github.com/user-attachments/assets/64f0f9da-d910-45ce-91cd-b6f7136db6c9)

**Album edit page**

![Image](https://github.com/user-attachments/assets/4d41f725-9fd4-4353-8801-c8e28b174843)

**Moderation page**

![Image](https://github.com/user-attachments/assets/703f6259-851a-46f9-b57d-1feef58720e3)

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
