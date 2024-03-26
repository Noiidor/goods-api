FROM golang:alpine3.19 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /goods-api ./cmd/api

FROM alpine:latest

COPY --from=build /goods-api /goods-api

CMD [ "/goods-api" ]