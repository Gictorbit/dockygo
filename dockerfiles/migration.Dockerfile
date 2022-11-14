# syntax=docker/dockerfile:1.4
ARG COMPRESS="true"
ARG VERSION="notset"
ARG BUILD_DATE="notset"
ARG HTTP_PROXY=""
ARG HTTPS_PROXY=""
ARG NO_PROXY=""
ARG MIGRATION_NAME="notset"

FROM migrate/migrate:latest AS migrate
# set proxies
ARG HTTP_PROXY
ARG HTTPS_PROXY
ARG NO_PROXY
ENV HTTP_PROXY=${HTTP_PROXY}
ENV HTTPS_PROXY=${HTTPS_PROXY}
ENV NO_PROXY=${NO_PROXY}

# install dependensies
RUN sed -i 's#dl-cdn.alpinelinux.org#alpine.global.ssl.fastly.net#g' /etc/apk/repositories
RUN apk add --update --no-cache upx

ARG COMPRESS
RUN if [ "$COMPRESS" = "true" ] ;then upx --best --lzma /usr/local/bin/migrate;fi

FROM scratch AS final
ARG MIGRATION_NAME
ARG VERSION
ARG BUILD_DATE

WORKDIR /sql
COPY "./sql" .
LABEL migration.build.version="${VERSION}"
LABEL migration.build.date="${BUILD_DATE}"
LABEL migration.name="${MIGRATION_NAME}"

WORKDIR /app
COPY --from=migrate /usr/local/bin/migrate .
ENTRYPOINT ["./migrate"]
CMD ["--help"]