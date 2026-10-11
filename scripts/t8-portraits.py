#!/usr/bin/env python3
"""Downloads the Tekken 8 character portraits bundled with the app from Wavu Wiki, as 256px
WebP files named by the character's name lowercased with only letters and digits kept
("Jack-8" -> jack8.webp).

The wiki's portraits are full-body renders. Most have the head near the top centre, so each
is cropped to a square from the top centre, a bit narrower than the image, to get a bust.

Usage: scripts/t8-portraits.py [output dir]
"""

import json
import os
import re
import subprocess
import sys
import tempfile
import urllib.parse
import urllib.request

API = "https://wavu.wiki/w/api.php"
SUFFIX = "_portrait.png"


def get(url):
    return urllib.request.urlopen(urllib.request.Request(url, headers={"User-Agent": "cfn-tracker portraits script"}))


def portraits():
    # The wiki has no category or title search for these, so go through every image.
    params = {"action": "query", "format": "json", "list": "allimages", "ailimit": "500", "aiprop": "url|size"}
    while True:
        data = json.load(get(f"{API}?{urllib.parse.urlencode(params)}"))
        for image in data["query"]["allimages"]:
            if image["name"].endswith(SUFFIX):
                yield image
        if "continue" not in data:
            return
        params.update(data["continue"])


def main():
    if len(sys.argv) > 2:
        sys.exit(__doc__.strip().splitlines()[-1])
    default = os.path.join(os.path.dirname(__file__), "..", "app", "gui", "src", "assets", "characters", "t8")
    out = sys.argv[1] if len(sys.argv) == 2 else default
    os.makedirs(out, exist_ok=True)
    with tempfile.TemporaryDirectory() as tmp:
        for image in portraits():
            character = image["name"].removesuffix(SUFFIX).replace("_", " ")
            name = re.sub(r"[^a-z0-9]", "", character.lower())
            png = os.path.join(tmp, f"{name}.png")
            with open(png, "wb") as f:
                f.write(get(image["url"]).read())
            side = min(image["width"], image["height"]) * 65 // 100
            x = (image["width"] - side) // 2
            webp = os.path.join(out, f"{name}.webp")
            subprocess.run(
                ["cwebp", "-quiet", "-crop", str(x), "0", str(side), str(side), "-resize", "256", "256", "-q", "80", png, "-o", webp],
                check=True,
            )
            print(webp)


if __name__ == "__main__":
    main()
