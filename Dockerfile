# syntax=docker/dockerfile:1
FROM golang:alpine

# Set destination for COPY
WORKDIR /app

RUN apk add --update alpine-sdk

# Download Go modules
COPY go.mod .
RUN go mod download

COPY . /app

# Build
RUN go build -buildvcs=false  -o . ./...


# second stage
FROM alpine:latest 
WORKDIR /root/

# install depends
COPY ./scripts/migrate ./
RUN apk add --update sqlite

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY --from=0 /app ./

# Run
CMD [ "./go_rest" ]