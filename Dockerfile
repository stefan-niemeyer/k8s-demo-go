# Build Stage: statisch bauen
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app

# Final Stage: minimal "FROM scratch"
FROM scratch
WORKDIR /
COPY --from=builder /src/app /app
COPY images/ images/
ENV PORT=3000 APP_VERSION=v1 APP_PICTURE=v1.jpg HOSTNAME=localhost
EXPOSE 3000
ENTRYPOINT ["/app"]
