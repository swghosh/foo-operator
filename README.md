# foo-operator
An example Kubernetes operator that syncs a deployment per FooManager resource, written using [reconciler.io/runtime](https://github.com/reconcilerio/runtime) with varifying simplified patterns over using just plain controller-runtime.

👀: [foo_reconciler.go](./internal/controller/foomanager_reconciler.go) called in [controller manager](./cmd/main.go) for details.
