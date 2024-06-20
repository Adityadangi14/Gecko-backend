FROM golang:1.22.4-alpine3.19
WORKDIR /server
COPY . /server
RUN go build /server
EXPOSE 3000
CMD [ "go","run","main.go" ]