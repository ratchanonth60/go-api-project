# Builder Stage
FROM golang:1.23-alpine AS builder

LABEL maintainer="By noobmaster"

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o apiserver cmd/runner.go 

# Nginx Stage
FROM nginx:alpine AS nginx

# Copy the Nginx config file into the container
COPY nginx.conf /etc/nginx/nginx.conf

# Copy the built Go API server from builder stage
COPY --from=builder /build/apiserver /usr/local/bin/

# Expose the port for the API and Nginx
EXPOSE 80
EXPOSE 8080

# Command to start both the API and Nginx
CMD ["sh", "-c", "/usr/local/bin/apiserver --config=env & nginx -g 'daemon off;'"]

