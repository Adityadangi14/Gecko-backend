FROM golang:1.22.4-alpine3.19
WORKDIR /server
COPY . /server
RUN go build /server
EXPOSE 3000
RUN sudo apt-get update && apt-get install -y supervisor
COPY supervisord.conf /etc/supervisor/conf.d/
CMD [ "go","run","main.go" ]