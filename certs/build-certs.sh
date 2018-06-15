#/bin/sh

SVC=$1

openssl genrsa -out ${SVC}/${SVC}.key 2048
openssl req -new -x509 -sha256 -key ${SVC}/${SVC}.key -out ${SVC}/${SVC}.crt -days 3650
openssl req -new -sha256 -key ${SVC}/${SVC}.key -out ${SVC}/${SVC}.csr
openssl x509 -req -sha256 -in ${SVC}/${SVC}.csr -signkey ${SVC}/${SVC}.key -out ${SVC}/${SVC}.crt -days 3650