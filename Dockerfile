FROM golang
MAINTAINER yuna

WORKDIR /code/golang_server/
ADD ./ ./
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' > /etc/timezone
ENV LANG C.UTF-8

# 依赖
ENV GOPROXY https://goproxy.cn

# 对外暴露的端口
EXPOSE 9001

RUN go build main.go -o t

ENTRYPOINT ["./t"]
