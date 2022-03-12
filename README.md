# RabbitMQ Cluster Kubernetes Operator

Forked from rabbitmq/cluster-operator to run on arm64 architecture.

```shell
# build modified image 
git clone git@github.com:fisruk/cluster-operator.git
docker build .

# push image to registry
docker tag [IMAGE-ID] registry.local/rabbitmq-cluster-operator:latest
docker push registry.local/rabbitmq-cluster-operator:latest

# install cluster operator
kubectl krew install rabbitmq
kubectl rabbitmq install-cluster-operator

kubectl set image -n rabbitmq-system deploy/rabbitmq-cluster-operator operator=registry.local/rabbitmq-cluster-operator
```
