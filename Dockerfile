FROM golang:1.26.3 as builder
LABEL maintainer "Sebastian Daehne <daehne@rshc.de>"
ENV GOOS=linux 
ENV GOARCH=386

RUN mkdir /build
WORKDIR /build
ADD . .
RUN go get -d -v ./... 
RUN go build -o dvb

FROM busybox
COPY --from=builder /build/dvb /dvb
CMD ["/dvb"]
