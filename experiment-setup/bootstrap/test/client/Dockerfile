FROM alpine:latest

WORKDIR /home/waect

COPY ./bootstrap/provisioning/client-provisioning.sh client-provisioning.sh
COPY ./client/scenario_helge.py scenario_helge.py

RUN chmod +x client-provisioning.sh

RUN ./client-provisioning.sh

CMD ["printenv"]