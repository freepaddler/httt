FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o httt

FROM scratch
LABEL authors="chu@vchu.net"
WORKDIR /app
COPY --from=builder /app/httt /app/

EXPOSE 8080

CMD ["/app/httt"]





