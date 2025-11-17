# Build Stage: statisch bauen
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app

# Final Stage: minimal
FROM alpine:latest
RUN apk add --no-cache bash curl

# download kubectl, make it executable and move it to a standard path
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
  && chmod +x kubectl \
  && mv kubectl /usr/local/bin/kubectl \

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
