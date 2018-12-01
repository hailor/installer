#!/bin/sh
mkdir -p {{pkiPath}}/etcd&&cd {{pkiPath}}/etcd
cfssl gencert -ca={{pkiPath}}/etcd/ca.crt -ca-key={{pkiPath}}/etcd/ca.key -config={{pkiPath}}/config/ca-config.json -profile=client {{pkiPath}}/config/etcd-client-csr.json | cfssljson -bare client&& mv client.pem client.crt && mv client-key.pem client.key
cfssl gencert -ca={{pkiPath}}/etcd/ca.crt -ca-key={{pkiPath}}/etcd/ca.key -config={{pkiPath}}/config/ca-config.json -profile=peer {{pkiPath}}/config/etcd-peer-csr.json | cfssljson -bare peer&& mv peer.pem peer.crt && mv peer-key.pem peer.key