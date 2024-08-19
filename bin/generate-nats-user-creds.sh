#!/bin/bash
#
# Description: Generates NATS User credentials
#
# Copyright (c) 2022 STY-Holdings Inc
# MIT License
# Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
# associated documentation files (the “Software”), to deal in the Software without restriction,
# including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the Software is furnished to
# do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all copies or
# substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
# BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
# NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
# DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#

# script variables
FILENAME=$(basename "$0")

# Private Variables

# shellcheck disable=SC2028
function print_usage() {
  echo "This will generate the user credentials for a NATS account."
  echo
  echo "Usage: $FILENAME -h | -a <NATS Account> -u <NATS Account User>"
  echo
  echo "Global flags:"
  echo "  -h                       Display help"
  echo "  -a <NATS account>        The name of the NATS account."
  echo "  -u <NATS account user>   The name of the NATS account user."
  echo
}

function set_variable() {
  cmd="${1}=$2"
  eval "$cmd"
}

function validate_arguments() {
  # shellcheck disable=SC2086
  if [ -z $NATS_ACCOUNT ]; then
    local Failed="true"
    echo "ERROR: You have to provide a NATS account."
  fi
  if [ -z "$NATS_ACCOUNT_USER" ]; then
    local Failed="true"
    echo "ERROR: The NATS account user was not provided."
  fi

  if [ "$Failed" == "true" ]; then
    print_usage
    exit 1
  fi
}

# Main function of this script
function run_script() {
  if [ "$#" == "0" ]; then
    echo "ERROR: No parameters where provided."
    print_usage
    exit 1
  fi

  while getopts 'ha:u:' OPT; do # see print_usage
    case "$OPT" in
    a)
      set_variable NATS_ACCOUNT "$OPTARG"
      ;;
    u)
      set_variable NATS_ACCOUNT_USER "$OPTARG"
      ;;
    h)
      print_usage
      exit 0
      ;;
    *)
      echo "ERROR: Please review the usage printed below:" >&2
      print_usage
      exit 1
      ;;
    esac
  done

  # Setup
  validate_arguments

  # Processing
  #
  # shellcheck disable=SC2034
  account=$2
  # shellcheck disable=SC2034
  user=$3

  echo " WARNING"
  echo " WARNING: You are creating a user credential file that has sensitive information!!!! "
  echo " WARNING           Handle with care and with system security in mind!!!!"
  echo " WARNING"
  echo " WARNING          The created file is located in $HOME/user-creds"
  echo

  # shellcheck disable=SC2086
  mkdir -p $HOME/user-creds
  # shellcheck disable=SC2086
  nsc generate creds --account $NATS_ACCOUNT --name $NATS_ACCOUNT_USER > $HOME/user-creds/$NATS_ACCOUNT_USER
  echo "Credentials have been generated."

  echo Done
}

run_script "$@"
