FROM gitpod/workspace-full

# More information: https://www.gitpod.io/docs/config-docker/
RUN sudo rm -rf /usr/bin/hd && \
    brew install linuxsuren/linuxsuren/hd && \
    hd install cli/cli
