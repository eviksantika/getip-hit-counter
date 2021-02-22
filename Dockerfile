FROM golang:1.15-buster AS build

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/apackeer/hit-counter

WORKDIR /go/src/github.com/apackeer/hit-counter
# Build go apps
RUN CGO_ENABLED=0 go build -o hit-counter main.go

FROM alpine:latest
# Export Listen port
EXPOSE 8080

# Copy file from build image
COPY --from=build /go/src/github.com/apackeer/hit-counter /opt
ENTRYPOINT /opt/hit-counter