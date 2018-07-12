FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go test -v ./...
RUN go install -v ./...

ENTRYPOINT ["app"]