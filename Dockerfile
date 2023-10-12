FROM alpine:3.18.4 as build
COPY --from=golang:1.21-alpine3.18 /usr/local/go/ /usr/local/go/
RUN GOBIN=/usr/local/bin/ /usr/local/go/bin/go install github.com/paololazzari/play@latest
RUN apk add --no-cache wget && \
    wget -q https://github.com/stedolan/jq/releases/download/jq-1.7/jq-linux64 && \
    mv jq-linux64 /usr/local/bin/jq && \
    wget -q https://github.com/mikefarah/yq/releases/download/v4.35.2/yq_linux_amd64 && \
    mv yq_linux_amd64 /usr/local/bin/yq

FROM alpine:3.18.4 as main
COPY --from=build /usr/local/bin/jq /usr/local/bin/jq
COPY --from=build /usr/local/bin/yq /usr/local/bin/yq
COPY --from=build /usr/local/bin/play /usr/local/bin/play
RUN apk add --no-cache bash=5.2.15-r5 && \
    rm -rf /var/cache/apk/* && \
    chmod +x /usr/local/bin/jq && \
    chmod +x /usr/local/bin/yq

COPY --from=build /usr/local/bin/play /usr/local/bin

ENTRYPOINT ["/usr/local/bin/play"]