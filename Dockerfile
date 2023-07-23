FROM golang

ENV GO111MODULE=on
ENV GOPROXY https://goproxy.io

ENV \
    PORT=9000 \
    HOST=0.0.0.0
EXPOSE 9000

RUN echo $GOPATH
WORKDIR $GOPATH/src
COPY . .
RUN go mod tidy

RUN go build -o main .
CMD ["./main", "&"]
