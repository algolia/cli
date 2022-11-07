FROM alpine
RUN apk update && apk upgrade && \
  apk add --no-cache ca-certificates
COPY algolia /bin/algolia
ENTRYPOINT ["/bin/algolia"]