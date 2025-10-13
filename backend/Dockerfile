FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd

FROM alpine:3.18
# create application directory
RUN mkdir -p /app /app/uploads
# copy the built binary into /app/app
COPY --from=builder /out/app /app/app
RUN chown root:root /app/uploads && chmod +x /app/app
EXPOSE 8080
ENTRYPOINT ["/app/app"]
