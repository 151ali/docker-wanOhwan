#FROM golang:1.11-alpine AS build
FROM golang:latest AS build
# Install tools you need to build your application
# Run `docker build --no-cache .` to update dependencies

RUN apt install git

# Install or update library dependencies
# List project dependencies with go.mod go.sum

COPY go.mod go.sum /go/src/server/
WORKDIR /go/src/server/

# Install library dependencies
RUN go mod download

# Copy the entire project and build it 
COPY . .
RUN go build -o /bin/server main.go

# This results in a single layer image

FROM ubuntu
COPY --from=build /bin/server /bin/server

EXPOSE 8008
CMD ["/bin/server"]