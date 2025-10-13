FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd

FROM alpine:3.18
COPY --from=builder /out/app /app
RUN mkdir -p /app/uploads && chown root:root /app/uploads
EXPOSE 8080
ENTRYPOINT ["/app"]
