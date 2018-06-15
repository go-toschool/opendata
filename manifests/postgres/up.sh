# got from: https://blog.oestrich.org/2015/08/running-postgres-inside-kubernetes/
kubectl create -f ./postgres-persistence.yml
kubectl create -f ./postgres-claim.yml
kubectl create -f ./postgres-svc.yml
kubectl create -f ./secret.yml
kubectl create -f ./postgres.yml
