FROM golang:1.19

RUN go env -w GO111MODULE=on

RUN go env -w GOPROXY=https://goproxy.cn,direct

MAINTAINER "xiaozuhui@tonatiuh.cn"

WORKDIR /home/workspace

ADD . /home/workspace

CMD go mod tidy

RUN go build main.go

EXPOSE 8088

ENTRYPOINT ["./main"]
