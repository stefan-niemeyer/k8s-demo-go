# Build Stage: statisch bauen
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app

# Final Stage: minimal
FROM alpine:latest
RUN apk add --no-cache bash curl

WORKDIR /
COPY --from=builder /src/app /app
COPY images/ images/

ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION

ARG APP_PICTURE
ENV APP_PICTURE=$APP_PICTURE

ARG UNSTABLE
ENV UNSTABLE=$UNSTABLE

ARG PORT=8080
ENV PORT=$PORT
EXPOSE $PORT

USER 65534

ENTRYPOINT ["/app"]
