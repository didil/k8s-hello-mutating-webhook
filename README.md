## K8s Hello Mutating Webhook
A Kubernetes Mutating Admission Webhook example, using Go.


#### Run tests
```
$ make test
```

#### Deploy
Define shell env:
```
$ # define env vars
$ export CONTAINER_REPO=quay.io/my-user/my-repo
$ export CONTAINER_VERSION=x.y.z
```

Build/Push Webhook 
```
$ make docker-build
$ make docker-push
```
* for this example you'll need to make the container repository public unless you'll be specifying ImagePullSecrets on the Pod

Deploy to K8s cluster
```
$ make k8s-deploy
```


#### Cleanup
Delete all k8s resources
```
$ make k8s-delete-all
```

