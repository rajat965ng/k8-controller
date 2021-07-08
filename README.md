[webpage]
- https://github.com/kubernetes/client-go 

[dependency]
- go get k8s.io/client-go/kubernetes@v0.20.0

[build]
-  env GOOS=linux GOARCH=arm go build

[role]
- kubectl  create role podAndDeploy --resource pods,deployments --verb list

[roleBindings]
- kubectl  create rolebinding podAndDeploy --role podAndDeploy --serviceaccount default:default