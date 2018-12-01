#!/bin/sh
kubeadm alpha phase certs all --config=/etc/kubernetes/kubeadm/master/kubeadm-master.yml
kubeadm alpha phase kubeconfig all --config=/etc/kubernetes/kubeadm/master/kubeadm-master.yml
kubeadm alpha phase controlplane all --config=/etc/kubernetes/kubeadm/master/kubeadm-master.yml
