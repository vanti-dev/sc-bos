# This dockerfile is for copying in pre-built binaries and assets.
# It is used by the GitHub Action docker-package.yml build process to create the release docker containers.

FROM alpine:3.18
LABEL vendor="Vanti Ltd"

# automatically populated by the builder
# will be 'linux/amd64' or 'linux/arm64' etc.
ARG TARGETPLATFORM

COPY .build/sc-bos/${TARGETPLATFORM} /app/
COPY .build/ops-ui/ /app/ops-ui/
COPY default/cfg/ /cfg/
COPY default/ui-config/ /app/ui-config/

EXPOSE 443
EXPOSE 23557
EXPOSE 7000-7999

VOLUME ["/cfg", "/data"]

ENTRYPOINT ["/app/sc-bos", "--appconf=/cfg/app.conf.json", "--sysconf=/cfg/system.conf.json", "--data=/data"]
