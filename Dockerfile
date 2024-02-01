FROM golang:latest

WORKDIR /go/src/organization-manager
COPY . .

EXPOSE 8080
RUN go build -v

CMD ["./organization-manager"]