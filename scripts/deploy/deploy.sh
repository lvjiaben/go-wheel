#!/bin/bash

# 设置变量
APP_NAME="go-admin"
VERSION="1.0.0"
NAMESPACE="default"

# 检查kubectl是否可用
if ! command -v kubectl &> /dev/null; then
    echo "kubectl could not be found"
    exit 1
fi

# 检查命名空间是否存在
if ! kubectl get namespace $NAMESPACE &> /dev/null; then
    echo "Creating namespace $NAMESPACE..."
    kubectl create namespace $NAMESPACE
fi

# 部署应用
echo "Deploying $APP_NAME..."
kubectl apply -f deployments/k8s/deployment.yaml -n $NAMESPACE

# 等待部署完成
echo "Waiting for deployment to complete..."
kubectl rollout status deployment/$APP_NAME -n $NAMESPACE

# 获取服务信息
echo "Service information:"
kubectl get service $APP_NAME -n $NAMESPACE

echo "Deployment completed successfully!" 