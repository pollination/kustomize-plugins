FROM golang:1.16.2-alpine3.13 AS build
RUN apk add --no-cache curl git openssh build-base
ENV GO111MODULE=on
COPY . /go/src/github.com/pollination/kustomize-plugins/
RUN cd /go/src/github.com/pollination/kustomize-plugins/ && make build-plugins

FROM alpine:3.13
RUN apk add --no-cache ca-certificates
ENV XDG_CONFIG_HOME=/
COPY --from=build /go/src/github.com/pollination/kustomize-plugins/dist/kustomize /kustomize
COPY --from=build /go/bin/kustomize /usr/local/bin/kustomize
ENTRYPOINT ["kustomize"]