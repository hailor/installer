#!/bin/sh
mkdir -p {{pkiPath}}&&cd {{pkiPath}}
cfssl gencert -initca {{pkiPath}}/config/k8s-front-proxy-ca-csr.json | cfssljson -bare front-proxy-ca && mv front-proxy-ca.pem front-proxy-ca.crt && mv front-proxy-ca-key.pem front-proxy-ca.key
cfssl gencert -ca={{pkiPath}}/front-proxy-ca.crt -ca-key={{pkiPath}}/front-proxy-ca.key -config={{pkiPath}}/config/ca-config.json -profile=client  {{pkiPath}}/config/k8s-front-proxy-client-csr.json | cfssljson -bare front-proxy-client && mv front-proxy-client-key.pem front-proxy-client.key && mv front-proxy-client.pem front-proxy-client.crt
