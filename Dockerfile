FROM golang:1.14-alpine as builder
RUN apk --no-cache add build-base git mercurial gcc
ADD . /minio-ups
WORKDIR /minio-ups
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./bin .

FROM scratch
COPY --from=builder /minio-ups/bin ./minio-ups
COPY .env .
ENTRYPOINT ["./minio-ups"]
EXPOSE 8080 