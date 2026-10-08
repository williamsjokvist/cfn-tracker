#!/usr/bin/env bash
# Simulates a network outage for CFN Tracker by blocking outbound traffic to the
# hosts it polls, so you can test the retry/recovery flow. Undo with
# scripts/unblock-network.sh.
#
# Blocks by IP rather than by program, since the SF6 tracker drives the system
# Chrome and a program rule would cut off all of your Chrome. Other sites that
# share these IPs (same CDN) will be blocked too while the rule is active.
#
# Windows: run from Git Bash as administrator. macOS: needs sudo.
set -euo pipefail

HOSTS=(www.streetfighter.com wank.wavu.wiki)
RULE_NAME="cfn-tracker-test-block"
PF_ANCHOR="com.apple/cfn-tracker-test"

resolve() {
  case "$(uname -s)" in
    MINGW* | MSYS* | CYGWIN*)
      powershell.exe -NoProfile -Command \
        "Resolve-DnsName -Name '$1' -ErrorAction SilentlyContinue | Where-Object { \$_.IPAddress } | ForEach-Object { \$_.IPAddress }" |
        tr -d '\r'
      ;;
    *)
      { dig +short A "$1"; dig +short AAAA "$1"; } | grep -E '^([0-9]+\.){3}[0-9]+$|:'
      ;;
  esac
}

ips=()
for host in "${HOSTS[@]}"; do
  resolved=$(resolve "$host" | sort -u)
  if [[ -z "$resolved" ]]; then
    echo "Could not resolve $host" >&2
    exit 1
  fi
  echo "$host -> $(echo $resolved)"
  while read -r ip; do ips+=("$ip"); done <<<"$resolved"
done

case "$(uname -s)" in
  MINGW* | MSYS* | CYGWIN*)
    if ! net session >/dev/null 2>&1; then
      echo "Run this from a Git Bash started as administrator." >&2
      exit 1
    fi
    netsh advfirewall firewall delete rule name="$RULE_NAME" >/dev/null 2>&1 || true
    remoteip=$(IFS=,; echo "${ips[*]}")
    MSYS_NO_PATHCONV=1 netsh advfirewall firewall add rule name="$RULE_NAME" \
      dir=out action=block remoteip="$remoteip" >/dev/null
    ;;
  Darwin)
    table=$(IFS=,; echo "${ips[*]}")
    printf 'table <cfn_tracker_test> { %s }\nblock drop out quick to <cfn_tracker_test>\n' "$table" |
      sudo pfctl -a "$PF_ANCHOR" -f - 2>/dev/null
    sudo pfctl -E >/dev/null 2>&1 || true
    # Drop existing connections so the block takes effect immediately.
    for ip in "${ips[@]}"; do sudo pfctl -k 0.0.0.0/0 -k "$ip" >/dev/null 2>&1 || true; done
    ;;
  *)
    echo "Unsupported OS: $(uname -s)" >&2
    exit 1
    ;;
esac

echo
echo "Blocked ${#ips[@]} addresses. Tracking should start retrying within one poll (~30s)."
echo "Run scripts/unblock-network.sh to restore."
