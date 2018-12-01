#!/bin/bash
docker run \
  -d \
  --privileged \
  --restart=on-failure:5 \
  --net=host \
  -v {{hostPkiPath}}:/etc/kubernetes/pki:ro \
  -v /var/lib/etcd:/var/lib/etcd:rw \
  --oom-kill-disable \
  --name=etcd \
  registry.cn-hangzhou.aliyuncs.com/choerodon-tools/etcd:v3.2.4 \
  /usr/local/bin/etcd \
  --data-dir=/var/lib/etcd \
  --name {{nodeName}} \
  --initial-advertise-peer-urls https://{{nodeIpAddress}}:2380 \
  --listen-peer-urls https://0.0.0.0:2380 \
  --advertise-client-urls https://{{nodeIpAddress}}:2379 \
  --listen-client-urls https://0.0.0.0:2379 \
  --initial-cluster {{etcdClusterPeers}} \
  --initial-cluster-state {{etcdState}} \
  --initial-cluster-token {{etcdToken}} \
  --client-cert-auth \
  --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt \
  --cert-file=/etc/kubernetes/pki/etcd/server.crt \
  --key-file=/etc/kubernetes/pki/etcd/server.key \
  --peer-client-cert-auth \
  --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt \
  --peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt \
  --peer-key-file=/etc/kubernetes/pki/etcd/peer.key