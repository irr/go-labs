#!/bin/bash
sudo rm -rf /opt/goapps /usr/local/bin/etcd* /usr/local/bin/nsq*
sudo mkdir -p /opt/goapps
sudo chown irocha: /opt/goapps
cd /opt/goapps
# https://github.com/coreos/etcd/releases/
curl -L  https://github.com/coreos/etcd/releases/download/v2.0.9/etcd-v2.0.9-linux-amd64.tar.gz -o etcd-v2.0.9-linux-amd64.tar.gz
tar xzvf etcd-v2.0.9-linux-amd64.tar.gz
rm -rf etcd; ln -s etcd-v2.0.9-linux-amd64 etcd
sudo ln -s /opt/goapps/etcd/etcd /usr/local/bin/etcd
sudo ln -s /opt/goapps/etcd/etcd /usr/local/bin/etcdctl
sudo ln -s /opt/goapps/etcd/etcd /usr/local/bin/etcd-migrate
# http://nsq.io/deployment/installing.html
wget https://s3.amazonaws.com/bitly-downloads/nsq/nsq-0.3.2.linux-amd64.go1.4.1.tar.gz
tar xfva nsq-0.3.2.linux-amd64.go1.4.1.tar.gz
rm -rf nsq;ln -s nsq-0.3.2.linux-amd64.go1.4.1 nsq
sudo ln -s /opt/goapps/
sudo ln -s /opt/goapps/nsq/bin/nsqadmin /usr/local/bin/nsqadmin
sudo ln -s /opt/goapps/nsq/bin/nsqd /usr/local/bin/nsqd
sudo ln -s /opt/goapps/nsq/bin/nsqlookupd /usr/local/bin/nsqlookupd
sudo ln -s /opt/goapps/nsq/bin/nsq_pubsub /usr/local/bin/nsq_pubsub
sudo ln -s /opt/goapps/nsq/bin/nsq_stat /usr/local/bin/nsq_stat
sudo ln -s /opt/goapps/nsq/bin/nsq_tail /usr/local/bin/nsq_tail
sudo ln -s /opt/goapps/nsq/bin/nsq_to_file /usr/local/bin/nsq_to_file
sudo ln -s /opt/goapps/nsq/bin/nsq_to_http /usr/local/bin/nsq_to_http
sudo ln -s /opt/goapps/nsq/bin/nsq_to_nsq /usr/local/bin/nsq_to_nsq
