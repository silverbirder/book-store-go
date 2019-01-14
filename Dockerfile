FROM golang:1.8
LABEL name "book-store-go"
WORKDIR /go/src/app
COPY . .
RUN go-wrapper download
RUN go-wrapper install
EXPOSE 3000
#ENTRYPOINT ["go-wrapper", "run"]