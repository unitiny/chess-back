FROM golang

RUN echo $GOPATH
RUN mkdir -p $GOPATH/src
WORKDIR $GOPATH/src
RUN cd $GOPATH/src
COPY . .

ENV GO111MODULE=on
ENV GOPROXY https://goproxy.io

ENV \
    PORT=9000 \
    HOST=0.0.0.0

EXPOSE 9000
RUN go run .