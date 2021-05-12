FROM gcr.io/distroless/static-debian10:latest

ARG VRESION

ARG TARGETARCH
COPY ./bin/cuem-operator_linux_${TARGETARCH}/cuem-operator /go/bin/cuem-operator

WORKDIR /go/bin

ENTRYPOINT ["/go/bin/cuem-operator"]
