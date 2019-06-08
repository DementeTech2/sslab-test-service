#!/bin/sh
DOMAIN="http://$(curl http://169.254.169.254/latest/meta-data/public-hostname --silent)"
jq -n --arg domain "$DOMAIN" '{ "apiEndpoint": $domain }' > config.json
aws s3 cp config.json s3://ssllabtest/config.json
rm config.json