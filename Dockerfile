# Stage 1: Build the Go binary
FROM --platform=$BUILDPLATFORM golang:1.20 as builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the source files
COPY . .

# Build the Go binary for the target platform
ARG TARGETPLATFORM
RUN CGO_ENABLED=0 GOOS=$(echo ${TARGETPLATFORM} | cut -d / -f1) GOARCH=$(echo ${TARGETPLATFORM} | cut -d / -f2) go build -a -o /app/main

# Stage 2: Create the final image
FROM busybox

# Copy the Go binary from the builder stage
COPY --from=builder /app/main /app/main

# Set the working directory
WORKDIR /app

# Expose the HTTP port
EXPOSE 8080

# Run the Go binary
CMD ["./main"]
