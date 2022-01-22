FROM golang:1.17-alpine As builder

RUN apk add build-base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o api .

FROM alpine:3.11.3

COPY --from=builder /app/api .
COPY --from=builder /app/.env .

EXPOSE 4001

CMD [ "/api" ]