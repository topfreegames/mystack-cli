mystack-cli
===========

The CLI for mystack

### Description
Mystack-cli is the Mystack interface to call the Mystack api. 

### Instalation
To build, run:
```
make build
```

Now the commands are available on ./bin/mysctl.

### Commands
#### Login
All the following commands require to the user to be logged in. 
```
./mysctl login 
```

It must be informed the router url with -s  and controller hostname with -o, for example:
```
./mysctl login -s http://my-router-url -o controller.mystack.com
```

It is possible to separate your config by environment (default is production):
```
./mysctl login -s http://my-router-url -o controller.mystack.com -e staging
```

#### Create config
First create a configuration file that specifies the apps that will run on Mystack cluster. 
For an example, check ./manifests/cluster.yaml.
To create it on Mystack, execute and inform the cluster name with -c and the path to file with -f:
```
./mysctl create config -e staging -c my-cluster -f /path/to/cluster-config
```

#### Delete config
To delete a config, inform its name with -c flag:
```
./mysctl delete config -e staging -c my-cluster
```

#### List config
To verify that your config was created, execute:
```
./mysctl get list -e staging
```

#### Get config
To read the saved cluster config, inform its name executing:
```
./mysctl get config -e staging -c my-cluster
```

#### Create cluster
After config creation, create the cluster itself. After creation, it returns the list of routes to the apps.
```
./mysctl create cluster -e staging -c my-cluster
```

To access these apps, either pass the route as Host header of every request, or use a DNS server to resolve the route to the Kubernetes IP. 

#### Get cluster
This method returns the app's routes. 
```
./mysctl get cluster -e staging -c my-cluster
```
