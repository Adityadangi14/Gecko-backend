FROM golang:1.22.4-alpine3.19
WORKDIR /server
COPY . /server
RUN go build /server
EXPOSE 3000
RUN  apk update && apk install -y supervisor
COPY supervisord.conf /etc/supervisor/conf.d/
CMD [ "go","run","main.go" ]