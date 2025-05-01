FROM golang:1.24

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN make build

EXPOSE 3000

CMD ["./bin/app"]