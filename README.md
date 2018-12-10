# RabbitMQ

Set up the RabbitMQ service in the Kubernetes cluster.

## Usage

### Up

The `up.sh` script performs the complete setup:

~~~bash
./up.sh
~~~

This script performs the following steps:

1. Deploy the resources in the `rabbitmq.yml` file (a deployment and a service) to the Kubernetes cluster
2. Replace the default RabbitMQ administrator user (username=*guest*, password=*guest*) with a new administrator user with a safe password
3. Save the credentials of the new RabbitMq administrator user in AWS Secrets Manager

### Down

You can use the `down.sh` script to delete all the resources created by `up.sh`:

~~~bash
./down.sh
~~~

## Requirements

To be able to execute the above scripts, the following constraints must be satisfied:

- Kubernetes cluster exists
- Kubernetes cluster is configured as the current context in the default *kubeconfig* file `~/.kube/config`
    - Executing `kubectl get nodes` should successfully access the right cluster
- `kubectl` is installed
- The AWS CLI is installed
- The AWS CLI is configured with the credentials and region that you want to use for creating resources
    - Check current credentials: `aws sts get-caller-identity`
    - Check current region: `aws configure get region`
