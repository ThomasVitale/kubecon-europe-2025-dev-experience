FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm 

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive && \
    apt-get -y install httpie && \
    rm -Rf /var/lib/apt/lists/*  && \
    rm -Rf /usr/share/doc && rm -Rf /usr/share/man  && \
    apt-get clean
