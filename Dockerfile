FROM golang:1.18

WORKDIR /app
COPY * ./
RUN go mod download

RUN go build ./cmd/arexAnalysis.go -o ./arexanalysis

EXPOSE 8090
CMD [ "./arexanalysis" ]