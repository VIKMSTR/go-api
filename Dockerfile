FROM busybox:glibc
RUN mkdir -p /db
COPY ./go-api /go-api
CMD ["/go-api", "--db-path", "/db/go-api.db"]

