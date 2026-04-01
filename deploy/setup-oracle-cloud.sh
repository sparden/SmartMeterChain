#!/bin/bash
# ============================================================
# SmartMeterChain — Oracle Cloud Always Free VM Setup
# Target: Ubuntu 22.04+ on Ampere A1 (4 OCPU, 24GB RAM)
# This script installs EVERYTHING and starts the full stack
# ============================================================
set -euo pipefail

echo "========================================="
echo " SmartMeterChain — Oracle Cloud Setup"
echo " Full Stack + Hyperledger Fabric"
echo "========================================="

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

step() { echo -e "\n${BLUE}[STEP]${NC} $1\n"; }
done_msg() { echo -e "${GREEN}[DONE]${NC} $1"; }

# ─── System Update ────────────────────────────────────────
step "1/10 — Updating system packages"
sudo apt update && sudo apt upgrade -y
sudo apt install -y \
    curl wget git unzip build-essential \
    apt-transport-https ca-certificates \
    software-properties-common gnupg lsb-release \
    jq tree htop tmux
done_msg "System updated"

# ─── Docker ───────────────────────────────────────────────
step "2/10 — Installing Docker"
if ! command -v docker &>/dev/null; then
    curl -fsSL https://get.docker.com | sudo sh
    sudo usermod -aG docker $USER
    sudo systemctl enable docker
    sudo systemctl start docker
    done_msg "Docker installed"
else
    done_msg "Docker already installed"
fi

# Docker Compose plugin
if ! docker compose version &>/dev/null; then
    sudo apt install -y docker-compose-plugin
fi
done_msg "Docker Compose ready"

# ─── Go 1.22 ─────────────────────────────────────────────
step "3/10 — Installing Go 1.22"
if ! command -v go &>/dev/null || [[ "$(go version)" != *"1.22"* ]]; then
    ARCH=$(dpkg --print-architecture)
    GO_TAR="go1.22.5.linux-${ARCH}.tar.gz"
    wget -q "https://go.dev/dl/${GO_TAR}"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "${GO_TAR}"
    rm "${GO_TAR}"

    # Add to PATH
    echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
    done_msg "Go $(go version) installed"
else
    done_msg "Go already installed"
fi

# ─── Node.js 20 LTS ──────────────────────────────────────
step "4/10 — Installing Node.js 20"
if ! command -v node &>/dev/null; then
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
    sudo apt install -y nodejs
    done_msg "Node $(node --version) installed"
else
    done_msg "Node already installed: $(node --version)"
fi

# ─── Rust ─────────────────────────────────────────────────
step "5/10 — Installing Rust"
if ! command -v rustc &>/dev/null; then
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
    source "$HOME/.cargo/env"
    echo 'source "$HOME/.cargo/env"' >> ~/.bashrc
    done_msg "Rust $(rustc --version) installed"
else
    done_msg "Rust already installed"
fi

# ─── Flutter ──────────────────────────────────────────────
step "6/10 — Installing Flutter (optional — skip on headless server)"
if ! command -v flutter &>/dev/null; then
    echo "  Skipping Flutter on server (use for mobile dev on local machine)"
    echo "  To install: snap install flutter --classic"
    done_msg "Flutter skipped (server environment)"
else
    done_msg "Flutter already installed"
fi

# ─── Hyperledger Fabric 2.5 ──────────────────────────────
step "7/10 — Installing Hyperledger Fabric 2.5 binaries"
if ! command -v peer &>/dev/null; then
    mkdir -p ~/fabric && cd ~/fabric
    curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh
    chmod +x install-fabric.sh
    ./install-fabric.sh --fabric-version 2.5.4 binary docker

    echo 'export PATH=$PATH:$HOME/fabric/bin' >> ~/.bashrc
    export PATH=$PATH:$HOME/fabric/bin
    cd ~
    done_msg "Fabric binaries installed"
else
    done_msg "Fabric already installed"
fi

# ─── Clone SmartMeterChain ────────────────────────────────
step "8/10 — Setting up SmartMeterChain project"
PROJECT_DIR="$HOME/SmartMeterChain"
if [ ! -d "$PROJECT_DIR" ]; then
    echo "  Project not found. Please clone or copy your project to $PROJECT_DIR"
    echo "  Example: git clone <your-repo-url> $PROJECT_DIR"
    mkdir -p "$PROJECT_DIR"
else
    done_msg "Project found at $PROJECT_DIR"
fi

# ─── Install Dependencies ────────────────────────────────
step "9/10 — Installing project dependencies"

# Backend Go deps
if [ -f "$PROJECT_DIR/backend/go.mod" ]; then
    cd "$PROJECT_DIR/backend"
    go mod tidy
    done_msg "Backend Go dependencies installed"
fi

# Frontend Node deps
if [ -f "$PROJECT_DIR/frontend/package.json" ]; then
    cd "$PROJECT_DIR/frontend"
    npm install
    done_msg "Frontend dependencies installed"
fi

# Simulator Rust deps
if [ -f "$PROJECT_DIR/simulator/Cargo.toml" ]; then
    cd "$PROJECT_DIR/simulator"
    cargo build --release 2>/dev/null || echo "  Simulator will build on first run"
    done_msg "Simulator compiled"
fi

cd ~

# ─── Firewall Rules ──────────────────────────────────────
step "10/10 — Configuring firewall (opening ports)"
sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 3000 -j ACCEPT  # Backend API
sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 3001 -j ACCEPT  # Frontend
sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 7050 -j ACCEPT  # Fabric Orderer
sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 7051 -j ACCEPT  # Fabric Peer
sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 5984 -j ACCEPT  # CouchDB
sudo netfilter-persistent save 2>/dev/null || true
done_msg "Firewall configured"

echo ""
echo "========================================="
echo " SETUP COMPLETE!"
echo "========================================="
echo ""
echo " To start the full stack:"
echo ""
echo " 1. Fabric Network:"
echo "    cd ~/SmartMeterChain/blockchain/network"
echo "    chmod +x scripts/*.sh"
echo "    ./scripts/setup-network.sh"
echo "    ./scripts/deploy-chaincode.sh"
echo ""
echo " 2. Backend API (new terminal/tmux):"
echo "    cd ~/SmartMeterChain/backend"
echo "    FABRIC_ENABLED=true go run main.go"
echo ""
echo " 3. Frontend (new terminal/tmux):"
echo "    cd ~/SmartMeterChain/frontend"
echo "    npm run dev -- --host 0.0.0.0"
echo ""
echo " 4. Simulator (new terminal/tmux):"
echo "    cd ~/SmartMeterChain/simulator"
echo "    cargo run --release -- --register"
echo ""
echo " Access:"
echo "   API:       http://<your-vm-ip>:3000"
echo "   Dashboard: http://<your-vm-ip>:3001"
echo ""
echo " TIP: Use tmux to run all 4 in one SSH session:"
echo "   tmux new -s smc"
echo "   Ctrl+B then % to split pane"
echo "========================================="
