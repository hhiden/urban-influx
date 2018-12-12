#!/bin/bash

function timeout() {
  SECONDS=0; TIMEOUT=$1; shift
  while eval $*; do
    sleep 5
    [[ $SECONDS -gt $TIMEOUT ]] && echo "ERROR: Timed out" && exit -1
  done
}

# Waits for all pods in the given namespace to complete successfully.
function wait_for_all_pods {
  echo "Waiting for all pods in " $1 " to be ready"
  timeout 300 "oc get pods -n $1 2>&1 | grep -v -E '(Running|Completed|STATUS)'"
}

# # Strimzi
oc adm policy add-cluster-role-to-user cluster-admin developer

oc apply -f https://github.com/strimzi/strimzi-kafka-operator/releases/download/0.8.2/strimzi-cluster-operator-0.8.2.yaml -n myproject
oc apply -f https://gist.githubusercontent.com/sjwoodman/ce869f8354e1e757f76c4ec8f9d06161/raw/25f60184f24d0d0fb841e0d8a6b17a9643072494/kafka-persistent.yaml -n myproject

wait_for_all_pods `oc project -q`

# KafkaEventSource CRD
oc create -f https://raw.githubusercontent.com/sjwoodman/eventing-sources/kafkaeventsource/contrib/kafka/kafkaeventsource-operator/deploy/crds/sources_v1alpha1_kafkaeventsource_crd.yaml

# Operator
oc create -f https://raw.githubusercontent.com/sjwoodman/eventing-sources/kafkaeventsource/contrib/kafka/kafkaeventsource-operator/deploy/service_account.yaml
oc create -f https://raw.githubusercontent.com/sjwoodman/eventing-sources/kafkaeventsource/contrib/kafka/kafkaeventsource-operator/deploy/role.yaml
oc create -f https://raw.githubusercontent.com/sjwoodman/eventing-sources/kafkaeventsource/contrib/kafka/kafkaeventsource-operator/deploy/role_binding.yaml
oc create -f https://raw.githubusercontent.com/sjwoodman/eventing-sources/kafkaeventsource/contrib/kafka/kafkaeventsource-operator/deploy/operator.yaml

#TODO: Wait for pods
wait_for_all_pods `oc project -q`
