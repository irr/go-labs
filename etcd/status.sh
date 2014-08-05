#!/bin/bash
PORT=${1:-4001}
echo "v1/machines"
http http://127.0.0.1:${PORT}/v1/machines
echo "v1/keys/_etcd/machines"
http http://127.0.0.1:${PORT}/v1/keys/_etcd/machines
echo "v1/leader"
http http://127.0.0.1:${PORT}/v1/leader