FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# set the time zone to Shanghai
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# add curl, used for http service checking, can be installed without it if deployed in k8s
RUN apk add curl

COPY configs/ /app/configs/
COPY caller /app/caller
RUN chmod +x /app/caller

# http port
EXPOSE 8080


WORKDIR /app

CMD ["./caller", "-c", "configs/caller.yml"]
# if you use the Configuration Center, caller.yml is changed to the Configuration Center configuration.
#CMD ["./caller", "-c", "configs/caller.yml", "-enable-cc"]
