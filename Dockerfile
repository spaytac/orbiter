FROM golang:1.8.3 as builder

WORKDIR /go/src/github.com/sarkk0x0/orbiter
ADD . /go/src/github.com/sarkk0x0/orbiter/
RUN make build

FROM scratch

COPY --from=builder /go/src/github.com/sarkk0x0/orbiter/bin/orbiter /bin/orbiter
# ENTRYPOINT ["orbiter"]
CMD ["/bin/orbiter", "daemon"]
