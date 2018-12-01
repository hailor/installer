#!/bin/sh
mkdir -p {{pkiPath}}/etcd&&cd {{pkiPath}}/etcd
cfssl gencert -ca={{pkiPath}}/etcd/ca.crt -ca-key={{pkiPath}}/etcd/ca.key -config={{pkiPath}}/config/ca-config.json -profile=server {{pkiPath}}/config/etcd-server-csr.json | cfssljson -bare server&& mv server.pem server.crt && mv server-key.pem server.key
