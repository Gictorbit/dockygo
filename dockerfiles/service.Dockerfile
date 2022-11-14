# syntax=docker/dockerfile:1.4
ARG GO_VERSION="1.19"
ARG GOPROXYURL="https://goproxy.io"
ARG COMPRESS="true"
ARG VERSION="notset"
ARG BUILD_DATE="notset"
ARG COMPANY_HOST="github.com/dolphinfms"
ARG GITHUB_TOKEN
ARG HTTP_PROXY=""
ARG HTTPS_PROXY=""
ARG NO_PROXY=""
ARG SERVICE_NAME=""

FROM golang:${GO_VERSION}-alpine AS builder
# set proxies
ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG NO_PROXY
ENV HTTP_PROXY=${HTTP_PROXY}
ENV HTTPS_PROXY=${HTTPS_PROXY}
ENV NO_PROXY=${NO_PROXY}
# install packages
RUN sed -i 's#dl-cdn.alpinelinux.org#alpine.global.ssl.fastly.net#g' /etc/apk/repositories
RUN apk --no-cache add --update ca-certificates tzdata upx git

# config git
ARG COMPANY_HOST
ARG GITHUB_TOKEN
ENV GOPRIVATE="${COMPANY_HOST}/*"
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@${COMPANY_HOST}".insteadOf "https://${COMPANY_HOST}"

# copy source code
WORKDIR /dolphinfms
COPY . .

# Get all of our dependencies
ARG GOPROXYURL
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod GOPROXY="${GOPROXYURL}" go mod download -x
# compile project
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w" -a -installsuffix cgo -buildvcs=false -o ./build/release/service ./cmd/...


ARG COMPRESS
RUN mkdir -p /final && \
    if [ "$COMPRESS" = "true" ] ;then upx --best --lzma -o /final/service ./build/release/service ;else cp ./build/release/service /final; fi


FROM scratch AS final
ARG VERSION
ARG BUILD_DATE
ARG SERVICE_NAME

WORKDIR /production
COPY --from=builder /final .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo


LABEL service.build.version="${VERSION}"
LABEL service.build.date="${BUILD_DATE}"
LABEL service.name="${SERVICE_NAME}"

ENTRYPOINT ["./service"]
