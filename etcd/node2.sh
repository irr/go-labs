#!/bin/bash
mkdir -p node2
etcd -s 127.0.0.1:7002 -c 127.0.0.1:8002 -C 127.0.0.1:7001 -d node2 -n node2