FROM alpine:3.13

RUN apk add ca-certificates
RUN update-ca-certificates

COPY kvconfig.yml /bin/kvconfig.yml
COPY bin/breakdown /bin/breakdown

CMD ["/bin/breakdown", "--addr=0.0.0.0:80"]
