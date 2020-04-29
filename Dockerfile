# build stage
FROM golang:1.14.1 as builder
## init setting
WORKDIR /workdir
ENV GO111MODULE="on"
## download packages
COPY go.mod go.sum ./
RUN go mod download
## build
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o minecraft-bot cmd/minecraft-bot/main.go

# run stage
FROM alpine:3.8 as app
WORKDIR /root/
## add Root certificates
RUN apk add --no-cache ca-certificates
## copy binary
COPY --from=builder /workdir/minecraft-bot .
## Run
ENTRYPOINT ["./minecraft-bot"]

