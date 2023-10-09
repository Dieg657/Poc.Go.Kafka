FROM golang:1.21-alpine AS BUILD

# Define WORKDIR
WORKDIR /app

# Copy MOD and SUM
COPY go.mod go.sum ./

# Download modules and verify modules
RUN go mod download && go mod verify

# Copy project files
COPY . .

# Build Go app to directory
RUN go build -o ./out/app .


# Define RUNTIME image
FROM alpine AS RUNTIME

# Define WORKDIR 
WORKDIR /app

# Install Certificates
RUN apk add ca-certificates

# Copy Go exec from BUILD
COPY --from=BUILD /app/out /app

# Expose port 8080
EXPOSE 8080

# Define entrypoint to Go built
CMD [ "./app" ]