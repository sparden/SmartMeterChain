#!/bin/bash
# ============================================================
# SmartMeterChain — Start All Services
# Works on any platform (local, cloud VM, Codespaces, Gitpod)
# ============================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

cleanup() {
    echo -e "\n${YELLOW}Stopping all services...${NC}"
    kill $BACKEND_PID $FRONTEND_PID $SIM_PID 2>/dev/null || true
    echo -e "${GREEN}All services stopped.${NC}"
    exit 0
}
trap cleanup SIGINT SIGTERM

echo "========================================="
echo " SmartMeterChain — Starting All Services"
echo "========================================="

# Check if Fabric is running
FABRIC_FLAG=""
if docker ps 2>/dev/null | grep -q "peer0.discom"; then
    echo -e "${GREEN}Fabric network detected — enabling blockchain mode${NC}"
    FABRIC_FLAG="FABRIC_ENABLED=true"
else
    echo -e "${YELLOW}Fabric not running — using mock mode${NC}"
    echo "  To start Fabric: cd blockchain/network && ./scripts/setup-network.sh"
fi
echo ""

# ─── Backend ─────────────────────────────────────────────
echo "Starting Backend API..."
cd "$PROJECT_DIR/backend"
if [ ! -d "vendor" ] && [ -f "go.mod" ]; then
    go mod tidy 2>/dev/null
fi
eval $FABRIC_FLAG go run main.go &
BACKEND_PID=$!
echo -e "  ${GREEN}Backend PID: $BACKEND_PID → http://localhost:3000${NC}"

# Wait for backend to be ready
echo "  Waiting for backend..."
for i in $(seq 1 15); do
    if curl -s http://localhost:3000/health >/dev/null 2>&1; then
        echo -e "  ${GREEN}Backend is up!${NC}"
        break
    fi
    sleep 1
done

# ─── Frontend ────────────────────────────────────────────
echo "Starting Frontend..."
cd "$PROJECT_DIR/frontend"
if [ ! -d "node_modules" ]; then
    npm install 2>/dev/null
fi
npm run dev -- --host 0.0.0.0 &
FRONTEND_PID=$!
echo -e "  ${GREEN}Frontend PID: $FRONTEND_PID → http://localhost:3001${NC}"

# ─── Simulator ───────────────────────────────────────────
echo "Starting Simulator..."
cd "$PROJECT_DIR/simulator"
sleep 3  # Let backend fully initialize
if [ -f "target/release/smartmeter-simulator" ]; then
    ./target/release/smartmeter-simulator --register &
else
    cargo run --release -- --register &
fi
SIM_PID=$!
echo -e "  ${GREEN}Simulator PID: $SIM_PID${NC}"

echo ""
echo "========================================="
echo -e " ${GREEN}ALL SERVICES RUNNING!${NC}"
echo "========================================="
echo ""
echo " Backend API:       http://localhost:3000"
echo " Frontend Dashboard: http://localhost:3001"
echo " Health Check:      http://localhost:3000/health"
echo ""
echo " Login Credentials:"
echo "   Admin:     admin / admin123"
echo "   Consumer:  consumer1 / consumer123"
echo "   Regulator: regulator / regulator123"
echo ""
echo " Press Ctrl+C to stop all services"
echo "========================================="

wait
