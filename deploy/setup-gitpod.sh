#!/bin/bash
# ============================================================
# SmartMeterChain — Gitpod Setup
# 50 hours/month free, Docker support, browser IDE
# Add this to your repo, then open: https://gitpod.io/#<repo-url>
# ============================================================
set -euo pipefail

echo "========================================="
echo " SmartMeterChain — Gitpod Setup"
echo "========================================="

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'
step() { echo -e "\n${BLUE}[STEP]${NC} $1\n"; }
done_msg() { echo -e "${GREEN}[DONE]${NC} $1"; }

# Gitpod has: Docker, Node.js, Python, Java
# Need to add: Go, Rust, Fabric

# ─── Go ───────────────────────────────────────────────────
step "1/5 — Installing Go 1.22"
if ! command -v go &>/dev/null; then
    wget -q "https://go.dev/dl/go1.22.5.linux-amd64.tar.gz"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
    rm go1.22.5.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
fi
done_msg "Go ready"

# ─── Rust ─────────────────────────────────────────────────
step "2/5 — Installing Rust"
if ! command -v rustc &>/dev/null; then
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
    source "$HOME/.cargo/env"
fi
done_msg "Rust ready"

# ─── Fabric ──────────────────────────────────────────────
step "3/5 — Installing Fabric 2.5"
mkdir -p ~/fabric && cd ~/fabric
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh
chmod +x install-fabric.sh
./install-fabric.sh --fabric-version 2.5.4 binary docker
export PATH=$PATH:$HOME/fabric/bin
echo 'export PATH=$PATH:$HOME/fabric/bin' >> ~/.bashrc
cd /workspace/SmartMeterChain 2>/dev/null || cd ~
done_msg "Fabric ready"

# ─── Dependencies ─────────────────────────────────────────
step "4/5 — Installing project dependencies"
PROJ="/workspace/SmartMeterChain"
[ ! -d "$PROJ" ] && PROJ="$GITPOD_REPO_ROOT"

if [ -f "$PROJ/backend/go.mod" ]; then
    cd "$PROJ/backend" && go mod tidy
    done_msg "Backend deps"
fi

if [ -f "$PROJ/frontend/package.json" ]; then
    cd "$PROJ/frontend" && npm install
    done_msg "Frontend deps"
fi

if [ -f "$PROJ/simulator/Cargo.toml" ]; then
    cd "$PROJ/simulator" && cargo build --release 2>/dev/null || true
    done_msg "Simulator built"
fi

# ─── Quick Start Script ──────────────────────────────────
step "5/5 — Creating start-all script"
cat > "$PROJ/start-all.sh" << 'STARTEOF'
#!/bin/bash
# Start all SmartMeterChain services in background
echo "Starting SmartMeterChain..."

# Backend
cd backend && go run main.go &
BACKEND_PID=$!
echo "Backend started (PID: $BACKEND_PID)"

# Wait for backend
sleep 3

# Frontend
cd ../frontend && npm run dev -- --host &
FRONTEND_PID=$!
echo "Frontend started (PID: $FRONTEND_PID)"

# Simulator
cd ../simulator && cargo run --release -- --register &
SIM_PID=$!
echo "Simulator started (PID: $SIM_PID)"

echo ""
echo "All services running!"
echo "  Backend:  http://localhost:3000"
echo "  Frontend: http://localhost:3001"
echo ""
echo "Press Ctrl+C to stop all"
wait
STARTEOF
chmod +x "$PROJ/start-all.sh"

echo ""
echo "========================================="
echo " GITPOD SETUP COMPLETE!"
echo "========================================="
echo ""
echo " Quick start (without Fabric):"
echo "   cd /workspace/SmartMeterChain"
echo "   ./start-all.sh"
echo ""
echo " With Fabric:"
echo "   cd blockchain/network && ./scripts/setup-network.sh"
echo "   Then ./start-all.sh in another terminal"
echo ""
echo " Gitpod auto-forwards all ports!"
echo "========================================="
