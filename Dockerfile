# syntax=docker/dockerfile:1

# specify the base image to  be used for the application, alpine or ubuntu
FROM registry.slauson.io/runner/go:latest

# create a working directory inside the image
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o iam-ms -buildvcs=false

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8080

# command to be used to execute when the image is used to start a container
CMD [ "./iam-ms" ]
