FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/commit-status-tracker \
    USER_UID=1001 \
    USER_NAME=commit-status-tracker

# install operator binary
COPY build/_output/bin/commit-status-tracker ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
