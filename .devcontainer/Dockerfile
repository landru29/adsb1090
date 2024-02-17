FROM mcr.microsoft.com/devcontainers/go:1.21-bullseye

ARG USERNAME=vscode

RUN apt update && apt install -y \
    vim \
    libusb-1.0-0-dev \
    git \
    build-essential \
    cmake \
    pkg-config \
    devscripts \
    curl \
    python3-pip \
    zsh

USER $USERNAME

RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.1.5/zsh-in-docker.sh)" -- \
    -t robbyrussell

RUN sudo mkdir -p /app/c && \
    sudo chown -R $USERNAME:$USERNAME /app && \
    git clone https://github.com/osmocom/rtl-sdr.git /app/c/rtl-sdr && \
    cd /app/c/rtl-sdr && sudo dpkg-buildpackage -us -uc && \
    sudo apt install -y /app/c/*.deb

RUN sudo update-alternatives --install /usr/bin/python python /usr/bin/python3.9 1 \
    && pip install pip --upgrade \
    && pip install pyModeS

RUN mknod /tmp/log p

COPY entrypoint.sh /
RUN sudo chmod ugoa+x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD [ ]