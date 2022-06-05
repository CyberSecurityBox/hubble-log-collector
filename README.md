# What is Hubble?

Hubble is a fully distributed networking and security observability platform
for cloud native workloads. It is built on top of [Cilium] and [eBPF] to enable
deep visibility into the communication and behavior of services as well as the
networking infrastructure in a completely transparent manner.


# What is Hubble-Log-Collector for?

Although Cilium recommends using `fluentd` to collect traffic logs but to make it simpler 
and more convenient to implement log collection and processing, I have added the ability 
to push logs to Kafka when using `hubble observe -o json` into the source code of Hubble.

## How does it work?

1. Let's create a namespace named `hubble-dev`
```
kubectl create namespace hubble-dev
```

2. Create a secret named `hubble-log-collector` to store connection information to Kafka
```
cat <<EOF > hubble-log-collector.yaml
apiVersion: v1
data:
  KAFKA_BROKERS: <base64>
  KAFKA_PASS: <base64>
  KAFKA_TOPIC: <base64>
  KAFKA_USER: <base64>
kind: Secret
metadata:
  name: hubble-log-collector
  namespace: hubble-dev
type: Opaque
EOF
```

3. Apply [secret](.k8s/secret.yaml) to cluster
```
kubectl apply -f hubble-log-collector.yaml
```

4. Create a Pod or Job to perform log collection

> Notes: \
> &nbsp; - See `hubble observe` support parameters at `hubble observe --help`. \
> &nbsp; - Recommend to use the parameters in `Filters Flags` to focus on logging specific workloads. \
> &nbsp; - Must have parameter `-f, --follow` Follow flows output. \
> &nbsp; - Must set setting environment variable `HUBBLE_SERVER`. By default, value is `hubble-relay.kube-system:80`. \
> &nbsp; - If the workload returns the error `K8S_CLUSTER_NAME is empty string`, the value must be set for the environment variable `K8S_CLUSTER_NAME`. By default, it is fetched from `http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-name`

- Apply [a pod](.k8s/pod.yaml)

  Used when you don't want to schedule log collection, 
  the user has to manually kill the pod.
  ```
  kubectl apply -f .k8s/pod.yaml
  ```
- Apply [a job](.k8s/job.yaml)

  Used to collect logs for a specific time period.
  ```
  kubectl apply -f .k8s/job.yaml
  ```