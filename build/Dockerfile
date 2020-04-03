FROM golang
LABEL "nebula-id"="captain"
WORKDIR /app
ADD dockercntrl /app/dockercntrl
ADD go.mod go.sum /app/
RUN go mod download
ADD . /app/
RUN go build -o main cmd/captain.go
ENTRYPOINT ["./main"]
