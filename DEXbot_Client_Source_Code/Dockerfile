#Get golang 1.17-buster as a base image
FROM golang:1.17-buster as builder

#Set build arguments
ARG BUILDOS
ARG BUILDARCH
ARG BUILDNAME



#Define the working directory in the container
WORKDIR /app

#Copy all files from root into the container
COPY . ./

#Use go mod tidy to handle dependencies
RUN go mod tidy


#Compile the binary
RUN env GOOS=$BUILDOS GOARCH=$BUILDARCH go build -o $BUILDNAME -trimpath -ldflags=-buildid=