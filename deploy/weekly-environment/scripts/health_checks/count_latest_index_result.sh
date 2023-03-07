#!/bin/bash

# This script is used to check if the latest elastic index has at least 200 results.
# It will make a request to the findings latest index and check if the response has at least 200 results.
# The script requires four arguments:
# 1. Elasticsearch URL
# 2. Kibana password
# 3. Request count
# 4. Request interval

ELASTICSEARCH_URL=$1
KIBANA_PASSWORD=$2
REQUEST_COUNT=$3
REQUEST_INTERVAL=$4
KIBANA_AUTH=elastic:${KIBANA_PASSWORD}

readonly MINIMAL_VALUE=200
readonly INDEX_NAME=logs-cloud_security_posture.findings_latest-default


for i in $(seq 1 "$REQUEST_COUNT"); do
  response="$(curl -X GET \
    --url "${ELASTICSEARCH_URL}/${INDEX_NAME}/_count" \
    -u "${KIBANA_AUTH}" \
    -H 'Cache-Control: no-cache' \
    -H 'Connection: keep-alive' \
    -H 'Content-Type: application/json' \
    -H 'kbn-xsrf: true')"
  count=$(echo "${response}" | jq -r '.count')
  echo "Request $i: $count results"
  if [ "$count" != "null" ] && [ "$count" -ge "$MINIMAL_VALUE" ]; then
    echo "The latest elastic index has at least $MINIMAL_VALUE results"
    exit 0
  fi
  sleep "$REQUEST_INTERVAL"
done

echo "The latest elastic index has less than $MINIMAL_VALUE results for $REQUEST_COUNT consecutive requests made within $REQUEST_INTERVAL"
exit 1