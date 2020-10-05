## K8s Hello Mutating Webhook
A Kubernetes Mutating Admission Webhook example, using Go.
This is a companion repository for the Article [Building a Kubernetes Mutating Admission Webhook: A “magic” way to inject a file into Pod Containers](https://medium.com/@didil/building-a-kubernetes-mutating-admission-webhook-7e48729523ed)

This is proof of concept code, make sure to review carefully before using in a production system.

#### Run tests
```
$ make test
```

#### Deploy
Define shell env:
```
# define env vars
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

#### Mutated pod example
```
$ k run busybox-1 --image=busybox  --restart=Never -l=hello=true -- sleep 3600
$ k exec busybox-1 -it -- ls /etc/config/hello.txt
# The output should be:
/etc/config/hello.txt
$ k exec busybox-1 -it -- sh -c "cat /etc/config/hello.txt"
# The output should be:
Hello from the admission controller !
```
We successfully mutated our pod spec and added an arbitary volume/file in there, yay !

#### Cleanup
Delete all k8s resources
```
$ make k8s-delete-all
```

