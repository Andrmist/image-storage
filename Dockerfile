FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go mod download && go build -o photo-storage .

FROM alpine

WORKDIR /app

COPY --from=builder /build/photo-storage /app/photo-storage

CMD ["./photo-storage"]
