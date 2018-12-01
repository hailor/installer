#!/bin/sh
mkdir -p {{pkiPath}}/etcd&&cd {{pkiPath}}/etcd
cfssl gencert -initca {{pkiPath}}/config/ca-csr.json | cfssljson -bare ca - && mv ca.pem ca.crt && mv ca-key.pem ca.key