#!/bin/bash
#
# STYH apt-update
#
# This will pull an apt upgrade list and generate commands
#
echo "Creating apt upgrade list"
touch /tmp/apt-upgrades.tmp
chmod 666 /tmp/apt-upgrades.tmp
# shellcheck disable=SC2034
apt list --upgradable --quiet > /tmp/apt-upgrades.tmp
echo "------"
echo "Generating apt upgrade commands"
/opt/utilities/apt-upgrades
rm /tmp/apt-upgrades.tmp
echo "done"
