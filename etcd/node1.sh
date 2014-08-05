#!/bin/bash
mkdir -p node1
etcd -s 127.0.0.1:7001 -c 127.0.0.1:8001 -d node1 -n node1