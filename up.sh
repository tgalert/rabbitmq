#!/bin/bash

set -e

# Coloured output
COL='\033[1;32m'
NOC='\033[0m'
echocln() {
  echo -e "$COL$*$NOC"
}
echoc() {
  echo -e -n "$COL$*$NOC"
}

#------------------------------------------------------------------------------#
# Deploy RabbitMQ
#------------------------------------------------------------------------------#

echocln "> Deploying RabbitMQ"
kubectl create -f rabbitmq.yml

#------------------------------------------------------------------------------#
# Wait for service endpoint
#------------------------------------------------------------------------------#

# Endpoint is either DNS name or IP address. This script works for both cases.
# JSONPath:
#   - https://kubernetes.io/docs/reference/kubectl/jsonpath/
# LoadBalancerIngress array:
#   - https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#loadbalanceringress-v1-core
echoc "> Waiting for service endpoint to be published"
while [[ -z "$endpoint" ]]; do
  sleep 1
  echoc .
  endpoint=$(kubectl get service rabbitmq -o jsonpath='{.status.loadBalancer.ingress[0].*}')
done
echo -e "\nSuccess: $endpoint"

#------------------------------------------------------------------------------#
# TODO: store service endpoint somewhere > AWS Systems Manager Parameter Store
#------------------------------------------------------------------------------#

#------------------------------------------------------------------------------#
# Wait for endpoint DNS resolution
#------------------------------------------------------------------------------#

echoc "> Waiting for endpoint DNS propagation (may take a few minutes)"
until host "$endpoint" >/dev/null; do
  sleep 1
  echoc .
done
echo -e "\nSuccess"

#------------------------------------------------------------------------------#
# Replace admin user
#------------------------------------------------------------------------------#

username=admin
password=$(cat /dev/urandom | LC_ALL=C tr -d -c 'a-zA-Z0-9' | head -c 32)
echocln "> Replacing default admin user {guest:guest} with {$username:**********}"

# Base URL of the RabbitMQ Management HTTP REST API
base_url="http://$endpoint:15672/api"

curl \
  -X PUT \
  -u guest:guest \
  -H "Content-Type: application/json" \
  -d "{\"password\":\"$password\",\"tags\":\"administrator\"}" \
  "$base_url/users/$username"

curl \
  -X DELETE \
  -u guest:guest \
  "$base_url/users/guest"

echo "Success"

#------------------------------------------------------------------------------#
# Store credentials of new admin user in AWS Secrets Manager
#------------------------------------------------------------------------------#

echocln "> Saving credentials of new admin user in AWS Secrets Manager"
env=prod

# Use the default AWS KMS customer master key (CMK) of the AWS account
# TODO: create a new KMS key just for this application
arn=$(aws secretsmanager create-secret \
  --name tgalert/$env/rabbitmq-admin \
  --description "Credentials of RabbitMQ administrator user." \
  --secret-string "{\"username\":\"$username\",\"password\":\"$password\"}" \
  --query ARN | tr -d \")

# To manually delete a secret without a recovery window, use:
# aws secretsmanager delete-secret --secret-id <NAME> --force-delete-without-recovery

echo "Success: $arn"
