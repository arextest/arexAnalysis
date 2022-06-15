FROM golang:1.18

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./

RUN go build cmd/arexAnalysis.go -o /arexAnalysis

EXPOSE 8090

CMD [ "/arexAnalysis" ]