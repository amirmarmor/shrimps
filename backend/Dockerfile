FROM gocv/opencv:4.5.3

ENV GOPATH /go

COPY . /go/src/gocv.io/x/gocv/

WORKDIR /app
COPY . ./
RUN GOARCH=arm go get
RUN GOARCH=arm go build -o /start
EXPOSE 1323

CMD ["/start"]

