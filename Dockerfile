FROM golang:1

ENV PROJECT="hmonit"

COPY . /${PROJECT}/
WORKDIR /${PROJECT}

RUN CGO_ENABLED=0 go build -mod=readonly -a -o /artifacts/${PROJECT}
COPY services.json /artifacts/services.json

FROM scratch
WORKDIR /
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /artifacts/* /

CMD ["/hmonit"]

