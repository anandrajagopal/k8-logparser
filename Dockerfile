FROM alpine:latest

WORKDIR /go

#COPY ./test.log .
COPY ./bin/main .

#RUN mkfifo test.log && chmod 777 * 

USER daemon

CMD [ "./main" ]