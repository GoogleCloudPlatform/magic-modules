FROM golang:1.20

WORKDIR /src
COPY . ./

RUN go build

ENTRYPOINT [ "/src/membership-checker" ]