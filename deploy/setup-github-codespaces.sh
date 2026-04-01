#!/bin/bash
# ============================================================
# SmartMeterChain — GitHub Codespaces Setup
# Pre-configured devcontainer with Docker-in-Docker support
# Run: Open repo in Codespaces (4-core, 16GB recommended)
# ============================================================
set -euo pipefail

echo "========================================="
echo " SmartMeterChain — Codespaces Setup"
echo "========================================="

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'
step() { echo -e "\n${BLUE}[STEP]${NC} $1\n"; }
done_msg() { echo -e "${GREEN}[DONE]${NC} $1"; }

# Codespaces already has: git, docker, node, python
# We need to add: Go, Rust, Fabric

# ─── Go 1.22 ─────────────────────────────────────────────
step "1/6 — Installing Go"
if ! command -v go &>/dev/null; then
    wget -q "https://go.dev/dl/go1.22.5.linux-amd64.tar.gz"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
    rm go1.22.5.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
fi
done_msg "Go $(go version)"

# ─── Rust ─────────────────────────────────────────────────
step "2/6 — Installing Rust"
if ! command -v rustc &>/dev/null; then
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
    source "$HOME/.cargo/env"
    echo 'source "$HOME/.cargo/env"' >> ~/.bashrc
fi
done_msg "Rust $(rustc --version)"

# ─── Hyperledger Fabric ──────────────────────────────────
step "3/6 — Installing Fabric 2.5 binaries + Docker images"
mkdir -p ~/fabric && cd ~/fabric
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh
chmod +x install-fabric.sh
./install-fabric.sh --fabric-version 2.5.4 binary docker
export PATH=$PATH:$HOME/fabric/bin
echo 'export PATH=$PATH:$HOME/fabric/bin' >> ~/.bashrc
cd /workspaces/SmartMeterChain 2>/dev/null || cd ~
done_msg "Fabric installed"

# ─── Backend Dependencies ────────────────────────────────
step "4/6 — Installing backend Go dependencies"
PROJ=$(find /workspaces -maxdepth 1 -name "SmartMeterChain" -type d 2>/dev/null || echo "")
if [ -z "$PROJ" ]; then
    PROJ="$HOME/SmartMeterChain"
fi

if [ -f "$PROJ/backend/go.mod" ]; then
    cd "$PROJ/backend" && go mod tidy
    done_msg "Backend deps ready"
fi

# ─── Frontend Dependencies ───────────────────────────────
step "5/6 — Installing frontend dependencies"
if [ -f "$PROJ/frontend/package.json" ]; then
    cd "$PROJ/frontend" && npm install
    done_msg "Frontend deps ready"
fi

# ─── Simulator Build ─────────────────────────────────────
step "6/6 — Building Rust simulator"
if [ -f "$PROJ/simulator/Cargo.toml" ]; then
    cd "$PROJ/simulator" && cargo build --release 2>/dev/null || echo "Will build on first run"
    done_msg "Simulator ready"
fi

echo ""
echo "========================================="
echo " CODESPACES SETUP COMPLETE!"
echo "========================================="
echo ""
echo " Ports will auto-forward. Open separate terminals:"
echo ""
echo " Terminal 1 — Fabric:"
echo "   cd blockchain/network && ./scripts/setup-network.sh"
echo ""
echo " Terminal 2 — Backend:"
echo "   cd backend && go run main.go"
echo ""
echo " Terminal 3 — Frontend:"
echo "   cd frontend && npm run dev"
echo ""
echo " Terminal 4 — Simulator:"
echo "   cd simulator && cargo run --release -- --register"
echo ""
echo " Codespaces auto-forwards ports 3000, 3001!"
echo "========================================="
