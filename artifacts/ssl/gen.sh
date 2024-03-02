#!/bin/bash

openssl genrsa 2048 > ca-key.pem
openssl req -x509 -nodes -days 36500 -key ca-key.pem -out ca-cert.pem -sha256 -subj "/CN=scheduler-plugins"

openssl req -newkey rsa:2048 -nodes -days 365000 -keyout server-key.pem -out server-req.pem -subj "/CN=queue-scheduler.kube-system.svc"
echo "subjectAltName = DNS:queue-scheduler.kube-system.svc" | openssl x509 -req -days 36500 -set_serial 01 -in server-req.pem -out server-cert.pem -CA ca-cert.pem -CAkey ca-key.pem -extfile -
