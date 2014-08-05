#!/bin/bash
mkdir -p node3
etcd -s 127.0.0.1:7003 -c 127.0.0.1:8003 -C 127.0.0.1:7001 -d node3 -n node3