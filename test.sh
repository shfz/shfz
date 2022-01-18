#!/bin/sh

cd `dirname $0`
set -eu

http --check-status --ignore-stdin --timeout=2.5 POST "http://localhost:53653/fuzz" name="test1" charset="abc" isGenetic:=false maxLen:=16 minLen:=8
http --check-status --ignore-stdin --timeout=2.5 POST "http://localhost:53653/fuzz" name="test2" charset="abcd" isGenetic:=false maxLen:=4 minLen:=2

curl -X POST http://localhost:53653/api -H "Content-Type: application/json" \
  -d '{"name": "login", "id": "0000", "fuzz": [{"name": "user", "text": "abcabc"}]}'

http --check-status --ignore-stdin --timeout=2.5 POST "http://localhost:53653/client/0000" \
  isClientError:=false

curl -X POST http://localhost:53653/server/0000 -H "Content-Type: application/json" \
  -d '{"isServerError": false, "frames": [{"name": "check", "file": "/main.go"}]}'

http --check-status --ignore-stdin --timeout=2.5 "http://localhost:53653/data"

curl -X POST http://localhost:53653/api -H "Content-Type: application/json" \
  -d '{"name": "login", "id": "1111", "fuzz": [{"name": "user", "text": "abcabc"}]}'

http --check-status --ignore-stdin --timeout=2.5 POST "http://localhost:53653/client/1111" \
  isClientError:=true \
  clientError="d is detected"

curl -X POST http://localhost:53653/server/1111 -H "Content-Type: application/json" \
  -d '{"isServerError": true, "serverError":"ZeroDivide", "serverErrorFile":"/main.go", "serverErrorLineNo": 10, "serverErrorFunc":"check1", "frames": [{"name": "check", "file": "/main.go"},{"name": "check1", "file": "/main.go"}]}'

http --check-status --ignore-stdin --timeout=2.5 "http://localhost:53653/data"

curl -F "hash=130d678cdc70fd679ce8e565ccdbac8c12f92098" -F "repo=shfz/shfz" -F "id=1587495459" -F "job=shfz" -F "number=5" -F "path=/app" http://localhost:53653/report
