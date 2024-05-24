#!/bin/bash
#
# STYH apt-update
#
# This will pull an apt upgrade list and generate commands
#
echo "Creating Directory"
ssh $IDENTITY $WORKING_AS@$INSTANCE_DNS_IPV4 "sudo mkdir /home/$SERVER_NAME/bin/"
echo "Copying script and binary"
scp $IDENTITY $ROOT_DIRECTORY/servers/$SERVER_NAME/bin/$SERVER_NAME $WORKING_AS@$INSTANCE_DNS_IPV4:/home/$SERVER_NAME/bin/apt-upgrades.sh
scp $IDENTITY $ROOT_DIRECTORY/servers/$SERVER_NAME/bin/$SERVER_NAME $WORKING_AS@$INSTANCE_DNS_IPV4:/home/$SERVER_NAME/bin/apt-upgrades
echo "done"
