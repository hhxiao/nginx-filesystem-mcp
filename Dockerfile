FROM alpine:3.22

RUN apk update && \
    apk upgrade && \
    apk add \
        bash \
        ca-certificates \
        wget \
    && rm -rf /var/cache/apk/*


COPY ./bin/mcp /opt/mcp

WORKDIR /opt

ENTRYPOINT ["/opt/mcp"]
