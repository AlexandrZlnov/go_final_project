FROM golang:latest AS builder

WORKDIR /go_final_project

COPY go.mod go.sum ./

RUN go mod download

COPY handlers ./handlers

COPY models ./models

COPY service ./service

COPY tests ./tests

COPY web ./web

COPY *.go *.md *.env ./

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

ENV TODO_PORT=7540 WEB_DIR="./web" TODO_DBFILE="scheduler.db"

RUN go build -o todo_app


FROM gcr.io/distroless/base-debian10

WORKDIR /go_final_project

COPY ./scheduler.db ./

COPY --from=builder /go_final_project /go_final_project

CMD ["./todo_app"]