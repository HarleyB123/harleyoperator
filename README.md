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
