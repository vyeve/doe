FROM vyeve/doe AS builder

FROM scratch

COPY --from=builder /usr/local/bin/client-ms /usr/local/bin/client-ms

WORKDIR /usr/local/bin/

CMD ["client-ms"]