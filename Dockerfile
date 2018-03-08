FROM golang:1.10.0 as builder
WORKDIR /go/src/github.com/ipedrazas/cluster-test
COPY . .
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
COPY --from=builder /go/src/github.com/ipedrazas/cluster-test/app /app

CMD [ "/app" ]
