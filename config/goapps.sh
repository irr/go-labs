#!/bin/bash
sudo rm -rf /opt/goapps /usr/local/bin/etcd* /usr/local/bin/nsq*
sudo mkdir -p /opt/goapps
sudo chown irocha: /opt/goapps
cd /opt/goapps
git clone https://github.com/pote/gpm.git && cd gpm
git checkout v1.3.2 # You can ignore this part if you want to install HEAD.
./configure
sudo make install
cd ..
git clone https://github.com/coreos/etcd.git
cd etcd && ./build
cd ..
sudo ln -s /opt/goapps/etcd/bin/etcd /usr/local/bin/etcd
sudo ln -s /opt/goapps/etcd/bin/etcdctl /usr/local/bin/etcdctl
git clone https://github.com/bitly/nsq.git
cd nsq
gpm install
go get -v github.com/bitly/nsq/internal/app
make
sudo ln -s /opt/goapps/nsq/build/apps/nsqadmin /usr/local/bin/nsqadmin
sudo ln -s /opt/goapps/nsq/build/apps/nsqd /usr/local/bin/nsqd
sudo ln -s /opt/goapps/nsq/build/apps/nsqlookupd /usr/local/bin/nsqlookupd
sudo ln -s /opt/goapps/nsq/build/apps/nsq_pubsub /usr/local/bin/nsq_pubsub
sudo ln -s /opt/goapps/nsq/build/apps/nsq_stat /usr/local/bin/nsq_stat
sudo ln -s /opt/goapps/nsq/build/apps/nsq_tail /usr/local/bin/nsq_tail
sudo ln -s /opt/goapps/nsq/build/apps/nsq_to_file /usr/local/bin/nsq_to_file
sudo ln -s /opt/goapps/nsq/build/apps/nsq_to_http /usr/local/bin/nsq_to_http
sudo ln -s /opt/goapps/nsq/build/apps/nsq_to_nsq /usr/local/bin/nsq_to_nsq
