FROM golang:1.18

WORKDIR /app
ADD . .
RUN go mod download
RUN go build -o arex-analysis ./cmd/main.go

EXPOSE 8090
CMD [ "./arex-analysis" ]