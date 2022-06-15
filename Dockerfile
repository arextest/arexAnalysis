FROM golang:1.18

WORKDIR /app
COPY . ./
RUN go mod download
RUN go build ./cmd/arexAnalysis.go

EXPOSE 8090
CMD [ "./cmd/arexAnalysis" ]