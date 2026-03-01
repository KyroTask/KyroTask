#!/bin/bash

# Suppress CGO Compiler Warnings (fixes the sqlite3-binding.c warning)
export CGO_CFLAGS="-w"

# Terminal Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

clear
echo -e "${PURPLE}"
echo "    ██╗  ██╗██╗   ██╗██████╗  ██████╗ ████████╗ █████╗  ███████╗██╗  ██╗"
echo "    ██║ ██╔╝╚██╗ ██╔╝██╔══██╗██╔═══██╗╚══██╔══╝██╔══██╗██╔════╝██║ ██╔╝"
echo "    █████╔╝  ╚████╔╝ ██████╔╝██║   ██║   ██║   ███████║███████╗█████╔╝ "
echo "    ██╔═██╗   ╚██╔╝  ██╔══██╗██║   ██║   ██║   ██╔══██║╚════██║██╔═██╗ "
echo "    ██║  ██╗   ██║   ██║  ██║╚██████╔╝   ██║   ██║  ██║███████║██║  ██╗"
echo "    ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝"
echo -e "${BLUE}                   🦊 F O X   M O D E 🦊${NC}"
echo ""
echo -e "${CYAN}🚀 Initializing KyroTask Ecosystem...${NC}"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  No .env file found. Creating from .env.example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✅ Created .env file${NC}"
    echo ""
    echo -e "${YELLOW}⚠️  IMPORTANT: Edit .env and set your TELEGRAM_BOT_TOKEN before continuing!${NC}"
    echo "   Get your bot token from @BotFather on Telegram"
    echo ""
    read -p "Press Enter after you've updated .env, or Ctrl+C to exit..."
fi

# Check if air is installed
if ! command -v air &> /dev/null; then
    echo -e "${YELLOW}⚠️  'air' is not installed. Installing...${NC}"
    go install github.com/air-verse/air@latest
    echo -e "${GREEN}✅ Installed air${NC}"
    echo ""
fi

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo -e "${CYAN}📦 Installing root dependencies...${NC}"
    npm install
    echo -e "${GREEN}✅ Root dependencies installed${NC}"
    echo ""
fi

if [ ! -d "mini-app/node_modules" ]; then
    echo -e "${CYAN}📦 Installing frontend dependencies...${NC}"
    cd mini-app && npm install && cd ..
    echo -e "${GREEN}✅ Frontend dependencies installed${NC}"
    echo ""
fi

# Download Go dependencies
echo -e "${CYAN}📦 Verifying Go dependencies...${NC}"
go mod download
echo -e "${GREEN}✅ Go environment ready${NC}"
echo ""

# Start development servers
echo -e "${PURPLE}🎯 IGNITING SERVERS...${NC}"
echo -e "   ${BLUE}Backend API:${NC}  http://localhost:3001"
echo -e "   ${GREEN}Web App UI:${NC}   http://localhost:5173"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all servers${NC}"
echo ""

# Launch the concurrent dev environment
npm run dev