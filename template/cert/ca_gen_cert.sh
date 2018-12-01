#!/bin/sh
mkdir -p {{pkiPath}}&&cd {{pkiPath}}
cfssl gencert -initca {{pkiPath}}/config/ca-csr.json | cfssljson -bare ca - && mv ca.pem ca.crt && mv ca-key.pem ca.key
