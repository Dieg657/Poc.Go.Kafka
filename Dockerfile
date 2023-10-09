FROM golang:1.21-alpine AS BUILD

RUN apk add alpine-sdk 
RUN apk --update add git

WORKDIR /app

# Copy MOD and SUM
COPY go.mod go.sum ./

# Download modules and verify modules
RUN go mod download && go mod verify

# Copy project files
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -tags musl -o ./out/app .

FROM alpine

WORKDIR /app

COPY --from=BUILD /app/out /app

CMD ["./app"]