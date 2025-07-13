#!/bin/bash

#minikube start

eval $(minikube docker-env)

./mockeventgenerator/build.sh
./leaderboard/build.sh
docker build -t frontend ./frontend

kubectl apply -f k8s/leaderboard-pvc.yaml
kubectl apply -f k8s/rabbitmq-deployment.yaml
kubectl apply -f k8s/leaderboard-deployment.yaml
kubectl apply -f k8s/mockeventgenerator-deployment.yaml
kubectl apply -f k8s/frontend-deployment.yaml
