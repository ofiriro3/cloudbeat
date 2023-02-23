#!/bin/bash

source ../utils.sh

# This script is used to create alerts for the vanilla integration.
# It will create a new slack connector and a new vanilla integration alerts.
# The script requires three arguments:
# 1. Kibana URL
# 2. Kibana password
# 3. Slack webhook

KIBANA_URL=$1
KIBANA_PASSWORD=$2
SLACK_WEBHOOK=$3
KIBANA_AUTH=elastic:${KIBANA_PASSWORD}

readonly SLACK_CONNECTOR_FILE=data/slack_connector.json
readonly VANILLA_ALERTS_FILE=data/vanilla_rules.ndjson

# Check if input is empty
if [ -z "$KIBANA_URL" ] || [ -z "$KIBANA_PASSWORD" ]; then
  echo "Kibana URL or Kibana password is empty"
  exit 1
fi

#######################################
# Creates new alerts for the vanilla integration
# Arguments:
#   $1: Kibana URL
#   $2: Kibana auth
#   $3: Alerts file
#   $4: Slack webhook url
# Returns:
#   None
#######################################
create_alerts_from_saved_object_file() {
  local kibana_url=$1
  local kibana_auth=$2
  local alerts_file=$3
  local slack_webhook_url=$4
  local slack_configuration_file=$5

  # Imports the slack connector and the rule alerts from a saved object file
  local import_alerts_response
  import_alerts_response=$(curl -X POST \
    "${kibana_url}/api/saved_objects/_import?overwrite=true" \
    -u "${kibana_auth}" \
    -H 'Cache-Control: no-cache' \
    -H 'Connection: keep-alive' \
    -H "kbn-xsrf: true" \
    --form file="@${alerts_file}")
  check_status_code_of_curl "${import_alerts_response}"

  # Get the connector id
  local connector_response
  connector_response=$(curl -X GET \
      "${kibana_url}/api/actions/connectors" \
      -u "${kibana_auth}" \
      -H 'Cache-Control: no-cache' \
      -H 'Connection: keep-alive' \
      -H "kbn-xsrf: true" \
      -H 'Content-Type: application/json')

  # Extracts the connector id of the slack connector
  local connector_id
  connector_id="$(echo "${connector_response}" |  jq '.[]  | select(.name == "#cloud-security-posture")' | jq -r '.id')"
  echo "Connector id: ${connector_id}"
  if [ -z "${connector_id}" ]; then
    echo "Connector id is empty"
    exit 1
  fi

  # Updates the slack connector with the webhook url
  local connector_configuration
  connector_configuration="$(jq --arg slack_webhook_url "${slack_webhook_url}" '.secrets.webhookUrl |= $slack_webhook_url' "${slack_configuration_file}")"
  echo "New connector configuration: ${connector_configuration}"
  local update_connector_response
  update_connector_response=$(curl -X PUT \
    "${kibana_url}/api/actions/connector/${connector_id}" \
    -u "${kibana_auth}" \
    -H 'Cache-Control: no-cache' \
    -H 'Connection: keep-alive' \
    -H "kbn-xsrf: true" \
    -H 'Content-Type: application/json' \
    -d "${connector_configuration}")

  check_status_code_of_curl "${update_connector_response}"

  # Get All rules
  local rules_response
  rules_response=$(curl -X GET \
      "${kibana_url}/api/alerts/_find" \
      -u "${kibana_auth}" \
      -H 'Cache-Control: no-cache' \
      -H 'Connection: keep-alive' \
      -H "kbn-xsrf: true" \
      -H 'Content-Type: application/json')

  # Covert to rules ids to array
  rules_ids=$(echo "${rules_response}" | jq -r '.data[].id')
  rules_ids=$(echo "${rules_ids}" | tr ' ' '\n' | jq -R . | jq -s .)

  echo "Rules ids: ${rules_ids}"

  # Enable all rules
  local enable_rule_response
  enable_rule_response=$(curl -X PATCH \
    "${kibana_url}/internal/alerting/rules/_bulk_enable" \
    -u "${kibana_auth}" \
    -H 'Cache-Control: no-cache' \
    -H 'Connection: keep-alive' \
    -H "kbn-xsrf: true" \
    -H 'Content-Type: application/json' \
    -d "{\"ids\": ${rules_ids}}")

  echo "Enable rule response: ${enable_rule_response}"
  check_status_code_of_curl "${enable_rule_response}"


#  for rule_id in ${rules_ids}; do
#    echo "Enabling rule: ${rule_id}"
#    local enable_rule_response
#    enable_rule_response=$(curl -X PATCH \
#      "${kibana_url}/internal/alerting/rules/_bulk_enable" \
#      -u "${kibana_auth}" \
#      -H 'Cache-Control: no-cache' \
#      -H 'Connection: keep-alive' \
#      -H "kbn-xsrf: true" \
#      -H 'Content-Type: application/json' \
#      -d "{\"ids\": ${rules_ids}}")
#
#    echo "Enable rule response: ${enable_rule_response}"
#    check_status_code_of_curl "${enable_rule_response}"
#  done
#  echo "Rules response: ${rules_response}"
  echo "Rules ids: ${rules_ids}"

}

# Create and enable alerts for the vanilla integration
create_alerts_from_saved_object_file "$KIBANA_URL" "$KIBANA_AUTH" "$VANILLA_ALERTS_FILE" "$SLACK_WEBHOOK" "$SLACK_CONNECTOR_FILE"
