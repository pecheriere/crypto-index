FROM golang:1.9

RUN mkdir -p /go/src/github.com/pecheriere/crypto-index
WORKDIR /go/src/github.com/pecheriere/crypto-index
COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]