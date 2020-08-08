# ******
FROM golang:latest AS build 
# ******

#FROM golang:1.12-alpine AS build
#RUN apk update && apk upgrade && \
#    apk add --no-cache bash git openssh

LABEL maintainer="IKHLEF Ali <ikhlef_ali@outlook.com>"

COPY go.mod go.sum /go/src/server/
WORKDIR /go/src/server/


RUN go mod download


COPY . .
RUN go build -o /bin/server app.go


# ******
FROM ubuntu
# ******

# FROM alpine
COPY --from=build /bin/server /bin/server

EXPOSE 8008
CMD ["/bin/server"]