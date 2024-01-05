FROM alpine:3.18
LABEL vendor="Vanti Ltd"

COPY .build/sc-bos /app/
COPY .build/ops-ui/ /app/ops-ui/
COPY static/cfg/ /app/cfg/
COPY static/ui-config/ /app/ui-config/

EXPOSE 443
EXPOSE 23557
EXPOSE 7000-7999

VOLUME ["/cfg", "/data"]

ENTRYPOINT ["/app/sc-bos", "--appconf=/cfg/app.conf.json", "--sysconf=/cfg/system.conf.json", "--data=/data"]
