FROM golang:1.23-alpine

RUN apk add --no-cache curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o metrics-processor

RUN chmod +x metrics-processor

CMD ["./metrics-processor"]