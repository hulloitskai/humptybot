##################################################
## BUILD STAGE
##################################################

FROM golang:alpine AS build

WORKDIR /go/src/github.com/steven-xie/humptybot
COPY . .

# Install dependencies...
RUN apk add curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ENV GOPATH="/go"
RUN dep ensure

# Create production binary...
RUN go build


##################################################
## PRODUCTION STAGE
##################################################

FROM alpine:3.7 as production
LABEL maintainer="Steven Xie <dev@stevenxie.me>"

RUN apk add ca-certificates

WORKDIR /app
COPY --from=build \
  /go/src/github.com/steven-xie/humptybot/humptybot \
  /go/src/github.com/steven-xie/humptybot/.env \
  ./

ENTRYPOINT [ "/app/humptybot" ]
