#!/usr/bin/env bash
# Downloads the Tekken 8 character portraits bundled with the app from Wavu Wiki, as 256px
# WebP files named by the character's name lowercased with only letters and digits kept
# ("Jack-8" -> jack8.webp).
#
# The wiki's portraits are full-body renders. Most have the head near the top centre, so each
# is cropped to a square from the top centre, a bit narrower than the image, to get a bust.
#
# Usage: scripts/t8-portraits.sh [output dir]
set -euo pipefail

out=${1:-"$(dirname "$0")/../app/gui/src/assets/characters/t8"}
mkdir -p "$out"

api=https://wavu.wiki/w/api.php
ua='cfn-tracker portraits script'
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# The wiki has no category or title search for these, so go through every image.
continue_from=''
while :; do
  page=$(curl -fsS -A "$ua" --get "$api" \
    --data-urlencode action=query --data-urlencode format=json --data-urlencode list=allimages \
    --data-urlencode ailimit=500 --data-urlencode aiprop='url|size' \
    ${continue_from:+--data-urlencode "aicontinue=$continue_from"})

  jq -r '.query.allimages[] | select(.name | endswith("_portrait.png")) | [.name, .url, .width, .height] | @tsv' <<<"$page" |
    while IFS=$'\t' read -r file url width height; do
      name=$(tr '[:upper:]' '[:lower:]' <<<"${file%_portrait.png}" | tr -cd 'a-z0-9')
      curl -fsS -A "$ua" -o "$tmp/$name.png" "$url"
      side=$(((width < height ? width : height) * 65 / 100))
      cwebp -quiet -crop $(((width - side) / 2)) 0 "$side" "$side" -resize 256 256 -q 80 \
        "$tmp/$name.png" -o "$out/$name.webp"
      echo "$out/$name.webp"
    done

  continue_from=$(jq -r '.continue.aicontinue // empty' <<<"$page")
  [ -n "$continue_from" ] || break
done
