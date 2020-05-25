FROM golang:1.14-alpine as builder

WORKDIR /go/src/redirector

COPY . .

RUN go get -d -v ./...
RUN go build -v ./...

FROM alpine:latest as runner

### Prepare user
RUN addgroup --gid 1001 redirector \
  && adduser \
  --home /dev/null \
  --gecos "" \
  --shell /bin/false \
  --ingroup redirector \
  --system \
  --disabled-password \
  --no-create-home \
  --uid 1001 \
  redirector

### copy binary
COPY --from=builder /go/src/redirector/redirector /usr/local/bin/redirector

ENTRYPOINT [ "redirector" ]
