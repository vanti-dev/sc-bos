#!/bin/sh
docker build -f Dockerfile-mac -t go-mac-builder .
docker run --rm -v $(pwd):/output go-mac-builder cp /app/mac-seeder /output/mac-seeder
./mac-seeder --appconf \
             ~/git/sc-bos/example/config/vanti-ugs/app.conf.json \
             --sysconf \
             ~/git/sc-bos/example/config/vanti-ugs/system.conf.json