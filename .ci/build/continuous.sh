#!/bin/bash

# Fail on any error.
set -e
# Display commands being run.
set -x

env | grep KOKORO

ruby --version
