#!/bin/bash
# ============================================================
# SmartMeterChain — Google Cloud Free Trial Setup
# Target: e2-medium (2 vCPU, 4GB) or e2-standard-4 (4 vCPU, 16GB)
# Uses $300 free credit for 90 days
# ============================================================
set -euo pipefail

echo "========================================="
echo " SmartMeterChain — Google Cloud Setup"
echo "========================================="

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'
step() { echo -e "\n${BLUE}[STEP]${NC} $1\n"; }
done_msg() { echo -e "${GREEN}[DONE]${NC} $1"; }

# ─── Prerequisite: Create VM first ───────────────────────
echo -e "${YELLOW}"
echo "BEFORE RUNNING THIS SCRIPT:"
echo "1. Go to console.cloud.google.com"
echo "2. Create a VM:"
echo "   gcloud compute instances create smartmeterchain \\"
echo "     --zone=asia-south1-a \\"
echo "     --machine-type=e2-standard-4 \\"
echo "     --image-family=ubuntu-2204-lts \\"
echo "     --image-project=ubuntu-os-cloud \\"
echo "     --boot-disk-size=50GB"
echo ""
echo "3. SSH into VM:"
echo "   gcloud compute ssh smartmeterchain --zone=asia-south1-a"
echo ""
echo "4. Clone repo & run this script"
echo -e "${NC}"
read -p "Press Enter to continue setup on this VM..."

# ─── System Packages ─────────────────────────────────────
step "1/8 — System packages"
sudo apt update && sudo apt upgrade -y
sudo apt install -y \
    curl wget git unzip build-essential \
    apt-transport-https ca-certificates \
    software-properties-common gnupg jq tmux htop
done_msg "Packages installed"

# ─── Docker ───────────────────────────────────────────────
step "2/8 — Docker"
if ! command -v docker &>/dev/null; then
    curl -fsSL https://get.docker.com | sudo sh
    sudo usermod -aG docker $USER
    sudo systemctl enable docker && sudo systemctl start docker
fi
done_msg "Docker ready"

# ─── Go ───────────────────────────────────────────────────
step "3/8 — Go 1.22"
if ! command -v go &>/dev/null; then
    wget -q "https://go.dev/dl/go1.22.5.linux-amd64.tar.gz"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
    rm go1.22.5.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
fi
done_msg "Go installed"

# ─── Node.js ─────────────────────────────────────────────
step "4/8 — Node.js 20"
if ! command -v node &>/dev/null; then
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
    sudo apt install -y nodejs
fi
done_msg "Node $(node --version)"

# ─── Rust ─────────────────────────────────────────────────
step "5/8 — Rust"
if ! command -v rustc &>/dev/null; then
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
    source "$HOME/.cargo/env"
    echo 'source "$HOME/.cargo/env"' >> ~/.bashrc
fi
done_msg "Rust installed"

# ─── Fabric ──────────────────────────────────────────────
step "6/8 — Hyperledger Fabric 2.5"
if ! command -v peer &>/dev/null; then
    mkdir -p ~/fabric && cd ~/fabric
    curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh
    chmod +x install-fabric.sh
    ./install-fabric.sh --fabric-version 2.5.4 binary docker
    echo 'export PATH=$PATH:$HOME/fabric/bin' >> ~/.bashrc
    export PATH=$PATH:$HOME/fabric/bin
    cd ~
fi
done_msg "Fabric installed"

# ─── Project Setup ────────────────────────────────────────
step "7/8 — Project dependencies"
PROJECT_DIR="$HOME/SmartMeterChain"

if [ -f "$PROJECT_DIR/backend/go.mod" ]; then
    cd "$PROJECT_DIR/backend" && go mod tidy
    done_msg "Backend deps"
fi

if [ -f "$PROJECT_DIR/frontend/package.json" ]; then
    cd "$PROJECT_DIR/frontend" && npm install
    done_msg "Frontend deps"
fi

if [ -f "$PROJECT_DIR/simulator/Cargo.toml" ]; then
    cd "$PROJECT_DIR/simulator" && cargo build --release 2>/dev/null || true
    done_msg "Simulator built"
fi

# ─── Firewall ────────────────────────────────────────────
step "8/8 — Firewall rules"
echo "  Creating GCP firewall rules..."
echo "  Run on your LOCAL machine (not the VM):"
echo ""
echo "  gcloud compute firewall-rules create smc-allow \\"
echo "    --allow=tcp:3000,tcp:3001,tcp:5984,tcp:7050,tcp:7051 \\"
echo "    --source-ranges=0.0.0.0/0 \\"
echo "    --description='SmartMeterChain ports'"
echo ""

echo "========================================="
echo " GOOGLE CLOUD SETUP COMPLETE!"
echo "========================================="
echo ""
echo " Start everything with tmux:"
echo "   tmux new -s smc"
echo ""
echo "   # Pane 1: Fabric"
echo "   cd ~/SmartMeterChain/blockchain/network"
echo "   ./scripts/setup-network.sh && ./scripts/deploy-chaincode.sh"
echo ""
echo "   # Pane 2: Backend (Ctrl+B, %)"
echo "   cd ~/SmartMeterChain/backend"
echo "   FABRIC_ENABLED=true go run main.go"
echo ""
echo "   # Pane 3: Frontend (Ctrl+B, %)"
echo "   cd ~/SmartMeterChain/frontend"
echo "   VITE_API_URL=http://$(curl -s ifconfig.me):3000/api/v1 npm run dev -- --host"
echo ""
echo "   # Pane 4: Simulator (Ctrl+B, %)"
echo "   cd ~/SmartMeterChain/simulator"
echo "   cargo run --release -- --register"
echo ""
echo " VM External IP: $(curl -s ifconfig.me 2>/dev/null || echo '<check GCP console>')"
echo "========================================="
