#!/bin/bash
set -e

minikube start --cpus=4 --memory=4096
minikube dashboard
