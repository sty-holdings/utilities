#!/bin/bash
#
# STYH apt-update
#
# This will pull an apt upgrade list and generate commands
#
echo "Generating upgrade list"
# shellcheck disable=SC2034
apt list --upgradable --quiet > /tmp/apt-upgrades.tmp
echo "------"
echo "Generating commands"
apt-update
