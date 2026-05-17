---
name: boost-factory-k8s
description: "Use when constructing a Kubernetes API client (clientset) via github.com/xgodev/boost/factory/contrib/k8s.io/client-go/v0. Covers NewClientset + variants and the typical use case (operators, controllers, jobs that talk to the cluster API). Triggers on imports under factory/contrib/k8s.io/client-go/, on questions about Kubernetes clientset construction in a boost service, or on in-cluster vs out-of-cluster auth."
license: MIT
metadata:
  author: jpfaria
  version: "0.1.0"
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

**REQUIRED BACKGROUND:** `boost-start`, `boost-wrapper-config`.

```go
import k8sfact "github.com/xgodev/boost/factory/contrib/k8s.io/client-go/v0"

cs := k8sfact.NewClientset(ctx)
// cs is *kubernetes.Clientset
```

Configure kubeconfig path, in-cluster vs out-of-cluster auth, QPS/burst under `boost.factory.k8s.*` (override `BOOST_FACTORY_K8S_*`).

## In-cluster vs out-of-cluster

In-cluster (running inside a k8s pod, picks up `/var/run/secrets/...`): leave `boost.factory.k8s.kubeconfig` empty.

Out-of-cluster (dev / one-off jobs): set `boost.factory.k8s.kubeconfig=~/.kube/config`.

## Red flags

| Red flag | Fix |
|---|---|
| `kubernetes.NewForConfig(restCfg)` directly with hand-built rest config | `k8sfact.NewClientset(ctx)` |
| Hardcoded kubeconfig path | `BOOST_FACTORY_K8S_KUBECONFIG` |
| Watch loops without `context.Context` cancellation | Pass the lifecycle context — leaked watches consume API server quota |
