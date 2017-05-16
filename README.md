mystack-cli
===========

The CLI for mystack

## About
Mystack-cli creates a cluster of services o Kubernetes. Create, delete, access them and read their logs using mystack-cli.

## Quickstart
* Download the latest [release](https://github.com/topfreegames/mystack-cli/releases)
* Login passing the controller's URL
```
mystack login http://controller.example.com
```
* See if there is an available config
```
mystack get configs
```
* If there is one named mycluster, for example, read the config details
```
mystack get config mycluster
```
* Create a cluster with that config
```
mystack create cluster mycluster
```
* Get the list of cluster's services and their URLs (all services are accessable at port 80)
```
mystack get cluster mycluster
```
* Read the service logs
```
mystack logs svc-name
```
* Follow the log stream with -f flag
```
mystack logs -f svc-name
```
* Bind local ports to services like dabatases
```
mystack port-forward mycluster
```
