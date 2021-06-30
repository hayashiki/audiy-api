# https://github.com/zenoss/gcloud-emulators/blob/e9a370d8cd678885a00cf2b0c64da9364acf4432/datastore/Dockerfile
FROM google/cloud-sdk:346.0.0-alpine
RUN apk --update --no-cache add openjdk8-jre
RUN gcloud components install beta cloud-datastore-emulator

ENV DATASTORE_CONSISTENCY=1.0 \
#    DATASTORE_DATA_DIR=/var/data \
    PORT=0.0.0.0:8081 \
    CLOUDSDK_CORE_PROJECT=audity-api-dev \
    DATASTORE_DISK_FLAG=--no-store-on-disk

EXPOSE 8081

CMD gcloud beta emulators datastore start \
    --host-port=$PORT \
    --no-store-on-disk \
    --consistency=$DATASTORE_CONSISTENCY
#    --data-dir=$DATASTORE_DATA_DIR
#     \
