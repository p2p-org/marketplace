FROM alpine:3.7

ENV ANSIBLE_VERSION 2.8.4

ENV BUILD_PACKAGES \
  bash \
  curl \
  tar \
  openssh-client \
  sshpass \
  git \
  python \
  py-boto \
  py-dateutil \
  py-httplib2 \
  py-jinja2 \
  py-paramiko \
  py-pip \
  py-yaml \
  docker \
  openssl \
  wget \
  ca-certificates
RUN wget -O /usr/local/bin/yq https://github.com/mikefarah/yq/releases/download/2.4.0/yq_linux_amd64 && chmod +x /usr/local/bin/yq
# If installing ansible@testing
#RUN \
#	echo "@testing http://nl.alpinelinux.org/alpine/edge/testing" >> #/etc/apk/repositories

RUN set -x && \
  \
  echo "==> Adding build-dependencies..."  && \
  apk --update add --virtual build-dependencies \
  gcc \
  musl-dev \
  libffi-dev \
  openssl-dev \
  python-dev && \
  \
  echo "==> Upgrading apk and system..."  && \
  apk update && apk upgrade && \
  \
  echo "==> Adding Python runtime..."  && \
  apk add --no-cache ${BUILD_PACKAGES} && \
  pip install --upgrade pip && \
  pip install python-keyczar docker-py && \
  \
  echo "==> Installing Ansible..."  && \
  pip install ansible==${ANSIBLE_VERSION} && \
  \
  pip install mitogen && \
  \
  echo "==> Cleaning up..."  && \
  apk del build-dependencies && \
  rm -rf /var/cache/apk/* && \
  \
  echo "==> Adding hosts for convenience..."  && \
  mkdir -p /etc/ansible /ansible && \
  echo "[local]" >> /etc/ansible/hosts && \
  echo "localhost" >> /etc/ansible/hosts

ENV ANSIBLE_STRATEGY_PLUGINS /usr/lib/python2.7/site-packages/ansible_mitogen/plugins/strategy
ENV ANSIBLE_STRATEGY mitogen_linear
ENV ANSIBLE_GATHERING smart
ENV ANSIBLE_HOST_KEY_CHECKING false
ENV ANSIBLE_RETRY_FILES_ENABLED false
# ENV ANSIBLE_ROLES_PATH /ansible/playbooks/roles
# ENV ANSIBLE_SSH_PIPELINING True
ENV PYTHONPATH /ansible/lib
ENV PATH /ansible/bin:$PATH
ENV ANSIBLE_LIBRARY /ansible/library

WORKDIR /work
