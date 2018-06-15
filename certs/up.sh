#/bin/sh

kubectl create -f castor/castor-secret.yml
kubectl create -f kanon/kanon-secret.yml
kubectl create -f saga/saga-secret.yml