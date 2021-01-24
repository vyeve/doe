FROM vyeve/doe AS builder

FROM scratch

COPY --from=builder /usr/local/bin/endpoint-ms /usr/local/bin/endpoint-ms

WORKDIR /usr/local/bin/

CMD ["endpoint-ms"]