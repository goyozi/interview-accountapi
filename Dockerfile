FROM golang:1.17

WORKDIR /usr/src/app

COPY go.* ./
RUN go mod download && go mod verify

COPY . .

CMD ["go", "test"]
