FROM alpine:3.14

RUN apk add --no-cache curl git make

ARG TARGETARCH
COPY ./bin/cuem_linux_${TARGETARCH}/cuem /bin/cuem

ENTRYPOINT /bin/cuem