FROM alpine:3.18.4 as build
COPY --from=golang:1.21-alpine3.18 /usr/local/go/ /usr/local/go/
RUN GOBIN=/usr/local/bin/ /usr/local/go/bin/go install github.com/paololazzari/play@latest

FROM alpine:3.18.4 as main
RUN apk add --no-cache bash=5.2.15-r5 && rm -rf /var/cache/apk/*
COPY --from=build /usr/local/bin/play /usr/local/bin

ENTRYPOINT ["/usr/local/bin/play"]