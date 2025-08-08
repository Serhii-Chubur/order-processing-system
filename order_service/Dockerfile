# -------------------------------------
# Dockerfile (development)
# -------------------------------------
FROM golang:1.24-alpine AS dev

# Робоча директорія
WORKDIR /usr/src/app

# 1) Скопіювати папку з цільовим мікросервісом (тут: `hello`)
COPY . product_service/

# Переходимо у директорію цільового сервісу
WORKDIR /usr/src/app/product_service

# 2) Завантажити залежності (використовуючи локальний go.mod сервісу)
RUN go mod download

# 3) Компілювати сервіс
CMD ["go", "run", "main.go"]