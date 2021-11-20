# lightweight container for go
FROM golang:alpine

# update container's packages and install git
RUN apk update && apk add --no-cache git bash
RUN apk add --update make
RUN apk add mysql-client

RUN mkdir /gocommerce

# set /app to be the active directory
WORKDIR /gocommerce

# copy all files from outside container, into the container
COPY . /gocommerce

RUN chmod 755 ./run.sh
RUN dos2unix ./run.sh

# download dependencies
RUN go mod tidy
RUN go get -u github.com/swaggo/swag/cmd/swag

# generate swagger docs
RUN swag init

ENTRYPOINT ./run.sh