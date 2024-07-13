# Телеграм бот для сбора достижений студентов


Телеграм бот, который добавляет и сохраняет индивидуальные достижения студентов. Преподаватель может выгружать и искать достижения студентов по всем возможным параметрам.
_____

# Как он работает?
Когда у пользователя запрашивается его роль (студент или преподаватель):

1.1. Если пользователь выбирает роль студента, он получает возможность добавлять, удалять и выполнять другие операции с достижениями.

1.2. Если пользователь выбирает роль преподавателя, он получает возможность просматривать достижения студентов, а также выгружать их в виде Excel-таблицы.

_____

# Как запустить свой собственный экземпляр бота?

Docker-compose: Dockerfile + PostgreSQL

### Что необходимо?

- Получите токен бота от [@BotFather](https://t.me/BotFather)
- Установить [Docker](https://www.docker.com/products/docker-desktop/)
- Поставить [image PostgreSQL](https://hub.docker.com/_/postgres)
____

# Инструкция
1. Клонировать репозиторий
```
git clone https://github.com/Ilusha-iz/tg-bot-student_achievements
```
2. Добавьте токен из BotFather  в переменную env в docker-compose.yml
```
version: '3.5'

services:
  db:
    image: postgres
    environment:
      POSTGRES_PASSWORD: <password>
    ports:
      - "5432:5432"

  bot:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      TOKEN: <your_telegram_bot_token>
      HOST: db
      USER: postgres
      PASSWORD: test
      DBNAME: postgres
    depends_on:
      - db

```

3. Создать контейнер
```
docker compose build
```

4. Запустить контейнер
```
docker compose up -d
```
____



