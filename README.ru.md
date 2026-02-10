[Русская версия](README.ru.md) | [English Version](README.md)

# 🎸 Tracker - сервис для прослушивания музыки

## 📃 Содержание

- [Скриншоты](#screenshots)
- [Идея проекта](#idea)
- [Основные фичи](#features)
- [Результат](#result)
- [Локальный запуск](#run)

<a name="screenshots"></a> 

## 🖼️ Скриншоты

**Домашняя страница**

![Image](https://github.com/user-attachments/assets/acf37b82-592b-4848-becf-9f470288e478)

**Регистрация**

![Image](https://github.com/user-attachments/assets/f12d4ade-f6de-4833-a017-51a6e2e8e049)

**Страница артиста**

![Image](https://github.com/user-attachments/assets/5c62f55c-787f-4a72-a071-8d4656b18a6c)

**Страница альбома**

![Image](https://github.com/user-attachments/assets/ae96145a-32bf-49a0-a640-f94d3654083e)

**Создание нового исполнителя**

![Image](https://github.com/user-attachments/assets/9dad0157-78c0-4d27-8c99-a720f1e07af8)

**Создание нового альбома**

![Image](https://github.com/user-attachments/assets/64f0f9da-d910-45ce-91cd-b6f7136db6c9)

**Редактирование альбома**

![Image](https://github.com/user-attachments/assets/4d41f725-9fd4-4353-8801-c8e28b174843)

**Админ-панель - модерация альбома**

![Image](https://github.com/user-attachments/assets/703f6259-851a-46f9-b57d-1feef58720e3)

<a name="idea"></a> 

## 💡 Идея

**Tracker** - альманах личных музыкальных предпочтений. Приложение для т.н. "музыкальных интровертов", где каждый может дополнять существующую музыкальную библиотеку

<a name="features"></a> 

## ✨ Основные фичи

- **Рекомендательная система**
  
  Все треки на главной странице, на странице любого исполнителя и в альбомах выбраны и отсортированы таким образом, чтобы соответствовать музыкальным предпочтениям конкретного пользователя

- **Нет границ между слушателем и исполнителем**
  
  Каждый может загружать свои песни, зарегистрировав профиль исполнителя. Пользователью предоставляется возможность видеть широкую статистику по прослушиваниям, аудитории и популярности его треков

<a name="result"></a> 

## ❗ Результат

На текущий момент полностью реализован функционал API.

Ознакомиться ендпоинтами можно либо в [SwaggerUI](localhost:8000/api/v1/docs), либо в файле [ENDPOINTS.md](https://github.com/illoprin/tracker/blob/master/ENDPOINTS.md)

<a name="run"></a> 

## ▶️ Локальный запуск

Запустить веб-приложение можно в режиме разработки (с hot-reload) через Docker

```bash
git clone https://github.com/illoprin/tracker.git
cd tracker
docker compose up -d --build
```

API будет доступно по адресу **localhost:8000**
