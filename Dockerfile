FROM alpine:3.22

LABEL maintainer="DevNet Engineering Team"
LABEL Description="DevNet PubHub Content MCP Server image" quay.expires-after="12w"

RUN apk update && \
    apk upgrade && \
    apk add \
        bash \
        ca-certificates \
        wget \
    && rm -rf /var/cache/apk/*


COPY ./bin/mcp /pubhub/mcp

WORKDIR /pubhub

ENTRYPOINT ["/pubhub/mcp"]
