FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN mkdir /data && touch /data/sqlite.db
RUN CGO_ENABLED=0 GOOS=linux go build -o /orbiter

FROM scratch

COPY --from=builder /orbiter /orbiter
COPY --from=builder /data/sqlite.db /data/sqlite.db
ENV DB_PATH=/data/sqlite.db

CMD ["/orbiter", "daemon"]
