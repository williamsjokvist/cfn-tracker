#!/usr/bin/env bash
# Downloads the SF6 character portraits bundled with the app from Buckler, as 256px WebP
# files named by the character's name lowercased with only letters and digits kept
# ("Chun-Li" -> chunli.webp).
#
# Usage: scripts/sf6-portraits.sh [output dir]
set -euo pipefail

out=${1:-"$(dirname "$0")/../app/gui/src/assets/characters/sf6"}
mkdir -p "$out"

# Capcom's site rejects requests without a browser user agent.
ua='Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15'
buckler=https://www.streetfighter.com/6/buckler/assets/images/material/character
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

slugs=$(curl -fsS -A "$ua" https://www.streetfighter.com/6/character | grep -oE 'character/[a-z_]+' | cut -d/ -f2 | sort -u)

for slug in $slugs; do
  # The character page and Buckler name a few characters differently.
  case $slug in
    gouki_akuma) image=gouki name=akuma ;;
    vega_mbison) image=vega name=mbison ;;
    ehonda) image=honda name=ehonda ;;
    *) image=$slug name=$slug ;;
  esac
  if ! curl -fsS -A "$ua" -o "$tmp/$name.png" "$buckler/character_${image}_r.png"; then
    echo "skipped $slug: no portrait at $buckler/character_${image}_r.png" >&2
    continue
  fi
  cwebp -quiet -resize 256 256 -q 80 "$tmp/$name.png" -o "$out/$name.webp"
  echo "$out/$name.webp"
done
