FROM golang:1.18 as builder
WORKDIR /go/src/app
RUN groupadd -g 10001 app \
  && useradd -u 10001 -g app app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip .
RUN ls -la .
RUN chmod 755 ./zoneinfo.zip
RUN make build

FROM scratch
USER app
ENV GOROOT /usr/local/go
COPY --from=builder /go/bin/app /go/bin/app
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/app/.conf /.conf
COPY --from=builder /go/src/app/.migrations /.migrations
COPY --from=builder /go/src/app/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
EXPOSE 8080
ENTRYPOINT ["/go/bin/app"]