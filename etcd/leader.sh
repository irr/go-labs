#!/bin/bash
PORT=${1:-4001}
echo "v1/leader"
http http://127.0.0.1:${PORT}/v1/leader