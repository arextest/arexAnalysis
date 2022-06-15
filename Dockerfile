FROM golang:1.18

WORKDIR /app
COPY * ./
RUN go mod download

RUN go build cmd/arexAnalysis.go -o /arexAnalysis

EXPOSE 8090
CMD [ "/arexAnalysis" ]