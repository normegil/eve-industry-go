#!/bin/sh
SERVER_NAME=$1
openstack server show -c addresses --format json $SERVER_NAME | jq -r '.addresses.\"Ext-Net\" | .[]' | grep "^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$"