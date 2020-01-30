FROM golang:alpine as builder

RUN apk update && apk add git 
COPY . $GOPATH/src/github.com/lacazethomas/goTodo
WORKDIR $GOPATH/src/github.com/lacazethomas/goTodo
RUN go get -d -v


RUN go build -o /go/bin/goToDoAPI


FROM alpine
EXPOSE 80
COPY --from=builder /go/bin/goToDoAPI /bin/goToDoAPI
ENTRYPOINT ["/bin/goToDoAPI"]