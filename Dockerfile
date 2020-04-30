FROM index.docker.io/gesquive/go-builder:latest AS builder

ENV APP=krypt

COPY dist/ /dist/
RUN copy-release

# =============================================================================
FROM gesquive/docker-base:latest
LABEL maintainer="Gus Esquivel <gesquive@gmail.com>"

# Import from builder
COPY --from=builder /app/${APP} /app/

ENTRYPOINT ["/app/krypt"]
