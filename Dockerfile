FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/authservice ./cmd/authservice

FROM alpine:3.22

RUN addgroup -S app && adduser -S -G app app

COPY --from=build /out/authservice /usr/local/bin/authservice

USER app

EXPOSE 50051

ENTRYPOINT ["authservice"]
