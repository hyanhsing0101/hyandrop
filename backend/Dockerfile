FROM golang:1.26-alpine AS dev

WORKDIR /app

RUN apk add --no-cache ca-certificates git tzdata

COPY go.mod go.sum* ./
RUN go mod download

CMD ["go", "run", "./cmd/server"]

FROM dev AS build

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/hyandrop ./cmd/server

FROM alpine:3.22 AS production

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=build /out/hyandrop /app/hyandrop

ENV HTTP_ADDR=:8080
EXPOSE 8080

CMD ["/app/hyandrop"]
