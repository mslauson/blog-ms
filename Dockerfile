# syntax=docker/dockerfile:1

# specify the base image to  be used for the application, alpine or ubuntu
FROM registry.slauson.io/runner/go:latest

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY . ./

# download Go modules and dependencies
RUN go mod download

# compile application
RUN go build -o post-ms -buildvcs=false

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8080

# command to be used to execute when the image is used to start a container
CMD [ "./post-ms" ]
