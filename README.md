# harleyoperator

```
operator-sdk init --domain harleyb123.com --repo github.com/HarleyB123/harleykubeoperator --skip-go-version-check
operator-sdk create api --kind SunriseSunset --resource --controller --version v1alpha1
make # To Test Controller
make generate
make manifests
make install run # For logs for Minikube
kubectl apply -f config/samples/_v1alpha1_sunrisesunset.yaml
```

TODO - Make it so that when the Spec Title changes, the site rebuilds with a different title. Also make it a service so its accessible outside the cluster.

Things for me to understand:

- The correct workflow when using the Make commands and what they do
- Fully understand the Reconciler code
- Understand the local dev workflow, currently building image, pushing image, minikube start, make deploy, kubectl apply and then check. Probably a faster method of iterating logic.
- Testing the Operator?


What I did wrong:

- Wrong domain when using the operator-sdk (hub.docker.com/repository/docker/harleyb123 or just hub.docker.com?). Sucks you can't change the domain after init
- Makefile needs to be changed manually (this sucks) for the BASE IMAGE and VERSION of the Controller
- Need to look into the RBAC Permissions and their manifests, I've now learnt that make manifests creates the RBAC for you!
- Need to push up the **controller** to docker hub or any other registry
