#!/usr/bin/env bash
set -euo pipefail

if [[ "${AXM_WEBHOOK_E2E:-}" != "true" ]]; then
  echo 'Set AXM_WEBHOOK_E2E=true to run the live webhook lifecycle.' >&2
  exit 1
fi

task_tmp="$(mktemp -d "${TMPDIR:-/tmp}/axm-webhook.XXXXXX")"
test_pid=''
cleanup() {
  if [[ -n "$test_pid" ]]; then
    kill -TERM "$test_pid" 2>/dev/null || true
    wait "$test_pid" 2>/dev/null || true
  fi
  rm -rf "$task_tmp"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

if [[ -z "${CLOUDFLARED_BIN:-}" ]]; then
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"
  case "$arch" in arm64|aarch64) arch=arm64 ;; x86_64) arch=amd64 ;; *) echo "Unsupported architecture: $arch" >&2; exit 1 ;; esac
  case "$os" in darwin) asset="cloudflared-darwin-${arch}.tgz" ;; linux) asset="cloudflared-linux-${arch}" ;; *) echo "Unsupported OS: $os" >&2; exit 1 ;; esac
  curl --fail --silent --show-error --location --connect-timeout 15 --max-time 120 "https://github.com/cloudflare/cloudflared/releases/latest/download/${asset}" --output "$task_tmp/$asset"
  if [[ "$os" == darwin ]]; then
    tar -xzf "$task_tmp/$asset" -C "$task_tmp"
  else
    mv "$task_tmp/$asset" "$task_tmp/cloudflared"
  fi
  chmod 700 "$task_tmp/cloudflared"
  export CLOUDFLARED_BIN="$task_tmp/cloudflared"
fi

repo_root="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$repo_root"
# Compile first so the live test binary, which owns the tunnel, directly receives signals.
go test -race -c -o "$task_tmp/webhook.test" ./axm/acceptance/webhook
"$task_tmp/webhook.test" -test.v -test.timeout=10m -test.run='^TestAcceptance_WebhookLifecycle$' &
test_pid=$!
wait "$test_pid"
test_pid=''
