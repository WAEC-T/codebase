FROM docker:27.3.1-dind-alpine3.20

WORKDIR /home/waect

COPY ./bootstrap/provisioning/server-provisioning.sh server-provisioning.sh

RUN chmod +x server-provisioning.sh

CMD ["dockerd-entrypoint.sh"]
