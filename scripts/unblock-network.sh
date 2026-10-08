#!/usr/bin/env bash
# Removes the block added by scripts/block-network.sh.
#
# Windows: run from Git Bash as administrator. macOS: needs sudo.
set -euo pipefail

RULE_NAME="cfn-tracker-test-block"
PF_ANCHOR="com.apple/cfn-tracker-test"

case "$(uname -s)" in
  MINGW* | MSYS* | CYGWIN*)
    if ! net session >/dev/null 2>&1; then
      echo "Run this from a Git Bash started as administrator." >&2
      exit 1
    fi
    if netsh advfirewall firewall delete rule name="$RULE_NAME" >/dev/null 2>&1; then
      echo "Unblocked."
    else
      echo "No block rule found, nothing to do."
    fi
    ;;
  Darwin)
    sudo pfctl -a "$PF_ANCHOR" -F all >/dev/null 2>&1 || true
    echo "Unblocked."
    ;;
  *)
    echo "Unsupported OS: $(uname -s)" >&2
    exit 1
    ;;
esac

echo "Tracking should recover on its next retry."
