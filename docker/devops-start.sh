#!/bin/bash

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

API_PORT=${4:-8000}
REDIS_PORT=${5:-6379}
PROMETHEUS_PORT=${6:-9090}
PUSHGATEWAY_PORT=${7:-9091}

if [ $# -lt 3 ]; then
    echo -e "${RED}Usage: $0 <version> <host> <web_port> [api_port] [redis_port] [prometheus_port] [pushgateway_port]${NC}"
    exit 1
fi

VERSION=$1
SERVER_HOST=$2
WEB_PORT=$3

for port in "$WEB_PORT" "$API_PORT" "$REDIS_PORT" "$PROMETHEUS_PORT" "$PUSHGATEWAY_PORT"; do
    if ! [[ "$port" =~ ^[0-9]+$ ]] || [ "$port" -lt 1 ] || [ "$port" -gt 65535 ]; then
        echo -e "${RED}Invalid port: $port${NC}"
        exit 1
    fi
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

detect_docker_compose() {
    if command -v docker-compose >/dev/null 2>&1; then
        echo "docker-compose"
    elif docker compose version >/dev/null 2>&1; then
        echo "docker compose"
    else
        echo ""
    fi
}

DOCKER_COMPOSE_CMD=$(detect_docker_compose)
if [ -z "$DOCKER_COMPOSE_CMD" ]; then
    echo -e "${RED}docker compose is required${NC}"
    exit 1
fi

if ! docker version >/dev/null 2>&1; then
    echo -e "${RED}Docker Engine is not healthy or not reachable.${NC}"
    echo "Please start or repair Docker Desktop / dockerd before running the phase1 stack."
    exit 1
fi

if [ ! -f "docker-compose.yml" ] || [ ! -f ".env" ]; then
    echo -e "${RED}docker-compose.yml and .env must exist in $(pwd)${NC}"
    exit 1
fi

set_env() {
    local key="$1"
    local value="$2"
    if grep -q "^${key}=" .env; then
        sed -i.bak "s|^${key}=.*|${key}=${value}|" .env
    else
        echo "${key}=${value}" >> .env
    fi
}

echo -e "${YELLOW}Updating docker environment...${NC}"
set_env "WEB_PORT" "$WEB_PORT"
set_env "API_PORT" "$API_PORT"
set_env "REDIS_PORT" "$REDIS_PORT"
set_env "PROMETHEUS_PORT" "$PROMETHEUS_PORT"
set_env "PUSHGATEWAY_PORT" "$PUSHGATEWAY_PORT"
set_env "DB_DIALECTS" "postgres"
set_env "IMAGE_HOST" "http://${SERVER_HOST}:${WEB_PORT}"
set_env "SERVER_PUBLIC_URL" "http://${SERVER_HOST}:${WEB_PORT}/"

echo -e "${YELLOW}Local source builds are enabled; version ${VERSION} will be used as an operator note only.${NC}"

echo -e "${YELLOW}Restarting services...${NC}"
$DOCKER_COMPOSE_CMD down 2>/dev/null || true
$DOCKER_COMPOSE_CMD up -d

sleep 10

SERVICES=("devops-redis" "devops-pushgateway" "devops-prometheus" "devops-api" "devops-web")
ALL_HEALTHY=true

for service in "${SERVICES[@]}"; do
    if docker ps --filter "name=$service" --filter "status=running" | grep -q "$service"; then
        echo -e "${GREEN}${service} is running${NC}"
    else
        echo -e "${RED}${service} is not running${NC}"
        ALL_HEALTHY=false
    fi
done

echo ""
if [ "$ALL_HEALTHY" = true ]; then
    echo -e "${GREEN}AutoOps phase1 stack is up${NC}"
    echo "  Web:         http://${SERVER_HOST}:${WEB_PORT}"
    echo "  API:         http://${SERVER_HOST}:${API_PORT}"
    echo "  PostgreSQL:  external (${DB_HOST:-set in .env}:${DB_PORT:-5432})"
    echo "  Redis:       ${SERVER_HOST}:${REDIS_PORT}"
    echo "  Prometheus:  http://${SERVER_HOST}:${PROMETHEUS_PORT}"
    echo "  Pushgateway: http://${SERVER_HOST}:${PUSHGATEWAY_PORT}"
else
    echo -e "${RED}Some services failed to start. Check logs with:${NC}"
    echo "  $DOCKER_COMPOSE_CMD logs -f"
    exit 1
fi
