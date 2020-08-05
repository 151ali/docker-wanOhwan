# parent image
FROM golang:latest
# set the working directory 
RUN mkdir /app
WORKDIR /app

# compile ...
COPY server.go .
RUN go build -o server .

# add metadata to the image to describe which port the container is listening on at runtime
EXPOSE 8080
# run the server
CMD ["./server"]

# copy the rest of app's source code from the host to the image file system.
COPY . .
