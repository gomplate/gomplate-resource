FROM golang:1.10-alpine AS build

ENV CGO_ENABLED 0
RUN mkdir -p /go/src/github.com/gomplate/gomplate-resource
WORKDIR /go/src/github.com/gomplate/gomplate-resource
COPY . /go/src/github.com/gomplate/gomplate-resource

RUN go build -o /out/in

FROM alpine:3.7

COPY --from=build /out/in /opt/resource/in
RUN ln -s /opt/resource/in /opt/resource/check

ARG CREATED
ARG REVISION
ARG VERSION

LABEL org.opencontainers.image.created=$CREATED \
      org.opencontainers.image.revision=$REVISION \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.name=gomplate-resource \
      org.opencontainers.image.authors=@hairyhenderson \
      org.opencontainers.image.source="https://github.com/gomplate/gomplate-resource"
