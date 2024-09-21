### k8s

Kubernetes是一个开源的容器编排平台，简称k8s。

换句话来说是管理容器的。既然k8s是管理容器的，那么docker产生的容器，或者不是Docker产生的容器k8s都可以管理。



基本概念



- Pod

  实例。一个Pod可以运行多个容器。比如你单个Pod可以运行MySQL、Redis等服务，一般不会这样做。

- Service

  逻辑上的服务，可以认为是你业务上某个服务的直接映射。

- Deployment

  用于管理Pod。

Pod和Service可以这样理解：假如我们有一个Web应用，它部署了3个实例。那么就是1个Web Service对应了3个Pod。







K8S调度的是容器，Docker容器运行的是各种镜像，要在k8s里面运行我们的webook，就需要把webook打包成一个镜像。

```
- 编译一个在Linux平台上执行的webook可执行文件
GOOS=linux GOARCH=arm go build -o webook .

- 运行docker，打包成一个镜像。
docker build -t flycash/webook:v0.0.1 .
```



打包成一个make docker命令：

```
rm webook || true
go mod tidy

GOOS=linux GOARCH=arm go build -o webook .
docker build -t flycash/webook:v0.0.1 .
```



powel shell

```
set GOARCH=amd64
set GOOS=linux
go build -o webook path/to/your/main.go
```



go env命令

```
- 查询go环境变量
go env

- 更改go环境变量
go env -w GOOS=linux

- arm64是ARM架构的CPU，amd64是X86架构的CPU，又叫X86_64
```



Makefile（make工具）

```
.PHONY docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f flycash/webook:v0.0.1
	@docker build -t flycash/webook:v0.0.1 .
```



### k8s配置

- deployment

  depolyment定义Pod的规格描述。

- service

  定义对外的可访问服务。

  一般会在service文件中要访问的deployment中定义好的Pod实例。



kubectl apply -f k8s-webook-deployment.yaml

kubectl get deployments



k8s是配置驱动的，通过apiVersion来解读配置。



在k8s里面，存储空间是被抽象为PersistentVolume（持久化卷）。因为k8s不知道容器运行的是什么，也不知道如何存储。因此便进行抽象，由具体的适配去实现。



Ingress

代表路由规则，前端的请求在经过ingress之后会转发到特定的Service上。和Service中的LoadBalancer相比，Service强调的是将流量转发到Pod上，二Ingress强盗的是发送到不同的Service上。

Ingress controller

一个Ingress controller可以控制住整个集群内部的所有Ingress（符合条件的Ingress）。



安装helm

使用helm安装ingress-nginx，运行：

```
helm upgrade --install ingress-nginx ingress-nginx --repo https://kubernetes.github.io/ingress-nginx --namespace ingress-nginx --create-namespace
```







修改代码后需要重新打包镜像

```
# 编译打包webook
set GOOS=linux
set GOARCH=amd64
go build -o webook .

# 构建webook service镜像
docker rmi -f flycash/webook:v0.0.1
docker build -t flycash/webook:v0.0.1 .






kubectl apply -f ./k8s-webook-deployment.yaml
kubectl apply -f ./k8s-redis-deployment.yaml
kubectl apply -f ./k8s-mysql-deployment.yaml

kubectl apply -f ./k8s-mysql-service.yaml
kubectl apply -f ./k8s-redis-service.yaml
kubectl apply -f ./k8s-webook-service.yaml

kubectl apply -f ./k8s-ingress-nginx.yaml
```



power shell

```
# 编译go程序
$Env:GOOS="linux"; $Env:GOARCH="amd64"; go build .

# 删除webook service镜像，# 重构镜像
docker rmi -f flycash/webook:v0.0.1
docker build -t flycash/webook:v0.0.1 .

kubectl delete deployment webook
kubectl delete service webook

kubectl apply -f ./k8s-webook-deployment.yaml
kubectl apply -f ./k8s-webook-service.yaml
```



go build编译标签

```
go build -tags=k8s .
```

