#!/bin/bash
SERVER_NAME=$1
IPV4_REGEX="^([0-9]{1,3}\.){3}[0-9]{1,3}$" 
ADDRESSES=($(openstack server show -c addresses --format json $SERVER_NAME | jq -r '.addresses."Ext-Net" | .[]'))
for address in "${ADDRESSES[@]}"
do
	if [[ $address =~ $IPV4_REGEX ]]; then
		echo $address
	fi
done
