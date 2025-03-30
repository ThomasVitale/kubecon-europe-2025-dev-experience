FROM mcr.microsoft.com/devcontainers/base:ubuntu-24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends httpie && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /usr/share/doc /usr/share/man && \
    apt-get clean

COPY ./install-sdkman.sh install-sdkman.sh
RUN chmod +x install-sdkman.sh && \
    ./install-sdkman.sh
