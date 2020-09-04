#!/bin/bash

function error_exit() {
  echo "$1" 1>&2
  exit 1
}

function check_deps() {
  test -f $(which jq) || error_exit "jq command not detected in path, please install it"
}

function parse_input() {
  # jq reads from stdin so we don't have to set up any inputs, but let's validate the outputs
  eval "$(jq -r '@sh "export REGION=\(.region) ASG_NAME=\(.asg_name)"')"
  if [[ -z "${REGION}" ]]; then export REGION=none; fi
  if [[ -z "${ASG_NAME}" ]]; then export ASG_NAME=none; fi
}

function get_instances_ids() {
  # get instances ids and put it on a JSON array
  RESULT=$(
    aws autoscaling describe-auto-scaling-groups \
      --auto-scaling-group-names ${ASG_NAME} \
      --region ${REGION} \
      --query AutoScalingGroups[].Instances[].InstanceId \
      --output json
  )
  jq -n --arg v "$RESULT" '{"instances_id": $v}'
}

# main()
check_deps
parse_input
get_instances_ids
