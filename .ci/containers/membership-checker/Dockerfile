FROM golang:1.18-alpine

WORKDIR /src
COPY . ./

RUN go build

ENTRYPOINT [ "/src/membership-checker" ]