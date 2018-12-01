#!/bin/sh
mkdir -p {{pkiPath}}&&cd {{pkiPath}}
cfssl gencert -ca={{pkiPath}}/ca.crt -ca-key={{pkiPath}}/ca.key -config={{pkiPath}}/config/ca-config.json -profile=server {{pkiPath}}/config/k8s-apiserver-csr.json | cfssljson -bare apiserver&& mv apiserver.pem apiserver.crt && mv apiserver-key.pem apiserver.key
cfssl gencert -ca={{pkiPath}}/ca.crt -ca-key={{pkiPath}}/ca.key -config={{pkiPath}}/config/ca-config.json -profile=client  {{pkiPath}}/config/k8s-apiserver-kubelet-client-csr.json | cfssljson -bare apiserver-kubelet-client  && mv apiserver-kubelet-client-key.pem apiserver-kubelet-client.key && mv apiserver-kubelet-client.pem apiserver-kubelet-client.crt
