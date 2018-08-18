##################################################
## BUILD STAGE
##################################################

FROM golang:alpine AS build

## Copy files
ENV GOPATH="/go"
WORKDIR ${GOPATH}/src/github.com/steven-xie/humptybot
COPY . .

## Install dependencies
RUN apk add curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure

## Create production binary
RUN go build


##################################################
## PRODUCTION STAGE
##################################################

FROM alpine:3.7 as production
LABEL maintainer="Steven Xie <dev@stevenxie.me>"

## Install dependencies
RUN apk add ca-certificates

## Copy files
WORKDIR /app
COPY --from=build \
  /go/src/github.com/steven-xie/humptybot/humptybot \
  /go/src/github.com/steven-xie/humptybot/.env \
  ./

## Set entrypoint
ENTRYPOINT [ "/app/humptybot" ]
