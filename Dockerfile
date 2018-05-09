FROM golang:1.7.4

MAINTAINER Duy Nguyen

COPY . ./app

WORKDIR ./app

RUN go install

ENV PORT 8080

EXPOSE 8080

ENTRYPOINT heroku-app