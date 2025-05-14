# spicedb-demo

## SpiceDB Cluster

```
minikube start
```

```
kubectl apply --server-side -f https://github.com/authzed/spicedb-operator/releases/latest/download/bundle.yaml
```

```
kubectl apply -f k8s/spicedb_cluster.yaml
```

```
kubectl port-forward deployment/dev-spicedb 50051:50051
```

## Create Relationship

```
brew install authzed/tap/zed
```

```
zed relationship create post:1 viewer user:emilia
zed relationship create post:2 viewer user:emilia

zed relationship create post:1 viewer user:bob
zed relationship create post:2 viewer user:bob
```

## Go Client

```
cd go
go run go/main.go
```

## TypeScript Client

```
cd ts
npm ci
npx tsx main.ts
```