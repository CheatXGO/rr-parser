FROM golang:latest
RUN go version
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go build -o main .
EXPOSE 8443
EXPOSE 5432