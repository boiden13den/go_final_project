FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler .

FROM ubuntu:latest

WORKDIR /app
COPY --from=builder /app/scheduler .
COPY web/ web/

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=12345

EXPOSE 7540

CMD ["./scheduler"]
