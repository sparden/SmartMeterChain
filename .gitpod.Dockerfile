FROM gitpod/workspace-full

# Go 1.22
RUN bash -c "source /home/gitpod/.sdkman/bin/sdkman-init.sh && \
    wget -q https://go.dev/dl/go1.22.5.linux-amd64.tar.gz && \
    sudo rm -rf /usr/local/go && \
    sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz && \
    rm go1.22.5.linux-amd64.tar.gz"
ENV PATH=$PATH:/usr/local/go/bin:/home/gitpod/go/bin

# Rust (already included in gitpod/workspace-full)

# Hyperledger Fabric binaries
RUN mkdir -p /home/gitpod/fabric && cd /home/gitpod/fabric && \
    curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && \
    chmod +x install-fabric.sh && \
    ./install-fabric.sh --fabric-version 2.5.4 binary
ENV PATH=$PATH:/home/gitpod/fabric/bin
