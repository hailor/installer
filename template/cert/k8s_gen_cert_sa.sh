#!/bin/sh
mkdir -p {{pkiPath}}&&cd {{pkiPath}}
openssl genrsa -out sa.key 1024  && openssl rsa -in sa.key -pubout -out sa.pub
