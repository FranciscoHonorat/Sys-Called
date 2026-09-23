#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
CLUSTER_NAME="${CLUSTER_NAME:-sys-called}"
NAMESPACE="${NAMESPACE:-sys-called}"
RELEASE="${RELEASE:-sys-called}"
TAG="${TAG:-dev-$(date +%Y%m%d%H%M%S)}"
CHART_DIR="$ROOT_DIR/infra/helm/sys-called"
KIND_CONFIG="$ROOT_DIR/infra/kind/cluster.yaml"
FRONTEND_URL="http://localhost:30080"

declare -A IMAGES=(
  [employees-service]="$ROOT_DIR/services/employees-service"
  [ticket-service]="$ROOT_DIR/services/ticket-service"
  [frontend]="$ROOT_DIR/frontend"
)

declare -A VALUE_KEYS=(
  [employees-service]="employeesService"
  [ticket-service]="ticketService"
  [frontend]="frontend"
)

log() {
  printf '\033[1;34m==>\033[0m %s\n' "$*"
}

fail() {
  printf '\033[1;31merro:\033[0m %s\n' "$*" >&2
  exit 1
}

require() {
  local missing=()
  for tool in "$@"; do
    command -v "$tool" >/dev/null 2>&1 || missing+=("$tool")
  done
  if ((${#missing[@]})); then
    fail "ferramentas não encontradas: ${missing[*]}"
  fi
}

cluster_exists() {
  kind get clusters 2>/dev/null | grep -qx "$CLUSTER_NAME"
}

ensure_cluster() {
  if cluster_exists; then
    log "cluster kind '$CLUSTER_NAME' já existe"
  else
    log "criando cluster kind '$CLUSTER_NAME'"
    kind create cluster --name "$CLUSTER_NAME" --config "$KIND_CONFIG" --wait 120s
  fi
  kubectl config use-context "kind-$CLUSTER_NAME" >/dev/null
}

build_and_load_images() {
  local name image output
  for name in "${!IMAGES[@]}"; do
    image="sys-called/$name:$TAG"
    log "build de $image"
    docker build -q -t "$image" "${IMAGES[$name]}" >/dev/null
    log "carregando $image no cluster"
    if ! output=$(kind load docker-image "$image" --name "$CLUSTER_NAME" 2>&1); then
      printf '%s\n' "$output" >&2
      fail "não foi possível carregar $image no cluster"
    fi
  done
}

deploy() {
  local args=()
  local name
  for name in "${!VALUE_KEYS[@]}"; do
    args+=(--set "images.${VALUE_KEYS[$name]}.tag=$TAG")
  done

  log "instalando o chart '$RELEASE' no namespace '$NAMESPACE'"
  helm upgrade --install "$RELEASE" "$CHART_DIR" \
    --namespace "$NAMESPACE" --create-namespace \
    "${args[@]}" \
    "$@" \
    --wait --timeout 10m
}

wait_for_frontend() {
  log "esperando o frontend responder em $FRONTEND_URL"
  for _ in $(seq 1 30); do
    if curl -sf "$FRONTEND_URL/" >/dev/null; then
      return 0
    fi
    sleep 2
  done
  fail "o frontend não respondeu em $FRONTEND_URL"
}

cmd_up() {
  require docker kind kubectl helm curl
  ensure_cluster
  build_and_load_images
  deploy "$@"
  wait_for_frontend
  log "pronto: $FRONTEND_URL (usuários de desenvolvimento com a senha senha123)"
}

cmd_deploy() {
  require docker kind kubectl helm
  cluster_exists || fail "o cluster '$CLUSTER_NAME' não existe; rode '$0 up' primeiro"
  kubectl config use-context "kind-$CLUSTER_NAME" >/dev/null
  build_and_load_images
  deploy "$@"
}

cmd_status() {
  require kubectl
  kubectl --context "kind-$CLUSTER_NAME" -n "$NAMESPACE" get pods,svc,pvc
}

cmd_logs() {
  require kubectl
  local target="${1:-ticket-service}"
  kubectl --context "kind-$CLUSTER_NAME" -n "$NAMESPACE" logs -f "deploy/$target" --all-containers
}

cmd_uninstall() {
  require helm kubectl
  log "removendo o release '$RELEASE' (os volumes e o secret são mantidos)"
  helm uninstall "$RELEASE" --namespace "$NAMESPACE" --kube-context "kind-$CLUSTER_NAME"
}

cmd_down() {
  require kind
  if cluster_exists; then
    log "apagando o cluster kind '$CLUSTER_NAME'"
    kind delete cluster --name "$CLUSTER_NAME"
  else
    log "cluster '$CLUSTER_NAME' não existe, nada a fazer"
  fi
}

usage() {
  cat <<USAGE
Uso: $(basename "$0") <comando> [argumentos extras do helm]

Comandos:
  up         cria o cluster kind (se preciso), faz o build das imagens e instala o chart
  deploy     refaz o build das imagens e atualiza o release num cluster existente
  status     mostra pods, services e volumes
  logs [app] acompanha os logs de um deployment (padrão: ticket-service)
  uninstall  remove o release, mantendo volumes e secret
  down       apaga o cluster kind inteiro

Variáveis: CLUSTER_NAME=$CLUSTER_NAME NAMESPACE=$NAMESPACE RELEASE=$RELEASE TAG=<gerada>
USAGE
}

main() {
  local command="${1:-}"
  shift || true
  case "$command" in
    up) cmd_up "$@" ;;
    deploy) cmd_deploy "$@" ;;
    status) cmd_status ;;
    logs) cmd_logs "$@" ;;
    uninstall) cmd_uninstall ;;
    down) cmd_down ;;
    -h | --help | help | "") usage ;;
    *) usage; fail "comando desconhecido: $command" ;;
  esac
}

main "$@"
