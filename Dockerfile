FROM golang:1.10.8 as builder
LABEL maintainer "Sebastian Daehne <daehne@rshc.de>"
ENV GOOS=linux 
ENV GOARCH=amd64

RUN mkdir /build
WORKDIR /build
ADD . .
RUN go get -d -v ./... 
RUN go build -o dvb

FROM alpine
COPY --from=builder /build/dvb /dvb
CMD ["/dvb"]
