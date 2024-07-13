# Шаг 1: Сборка
FROM golang:1.22.4-alpine AS build

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY code/go.mod code/go.sum ./
RUN go mod download

# Копируем исходный код
COPY code/ .

# Сборка приложения
RUN go build -o code .

# Шаг 2: Создание конечного образа
FROM alpine

ENV LANGUAGE="en"
WORKDIR /root/

# Копируем исполняемый файл из предыдущего этапа
COPY --from=build /app/code .

RUN apk add --no-cache ca-certificates

EXPOSE 80/tcp

CMD ["./code"]
