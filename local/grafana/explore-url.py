import json
import pathlib
import shutil
import subprocess
import sys
import urllib.parse

GRAFANA = "http://localhost:3000"
DATASOURCE = {"type": "loki", "uid": "micro-loki"}

query = (pathlib.Path(__file__).with_name("explore-query.logql")
         .read_text(encoding="utf-8").strip())

panes = {
    "logs": {
        "datasource": DATASOURCE["uid"],
        "queries": [{"refId": "A", "datasource": DATASOURCE, "expr": query,
                     "queryType": "range"}],
        "range": {"from": "now-1h", "to": "now"},
    }
}

params = urllib.parse.urlencode(
    {"schemaVersion": "1", "orgId": "1", "panes": json.dumps(panes, ensure_ascii=False)}
)
url = f"{GRAFANA}/explore?{params}"

if "--open" in sys.argv:
    opener = {"darwin": "open", "linux": "xdg-open"}.get(sys.platform)
    if opener and shutil.which(opener):
        subprocess.run([opener, url], check=False)
        print("Explore открыт в браузере")
        raise SystemExit

print(url)
