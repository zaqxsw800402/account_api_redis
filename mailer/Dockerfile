FROM golang:1.17-alpine as builder

RUN apk add build-base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o mailer .

FROM alpine:3.11.3

COPY --from=builder /app/mailer .
COPY --from=builder /app/.env .

EXPOSE 4002

CMD [ "/mailer" ]