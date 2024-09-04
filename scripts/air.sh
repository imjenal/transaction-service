#!/bin/sh
# This script is used to run air for hot reload. Don't run this script directly. Use the command `make dev` instead.

# In macOS, the default limit is 256 open files. This is not enough for air (https://github.com/air-verse/air) keep track of all the file changes for hot reload.
# This script will increase the limit to 20000 temporarily for the current session.
if [[ "${OSTYPE}" == *"darwin"* ]]; then
  echo "macOS detected. Increasing ulimit to 20000"
  ulimit -n 20000
fi

# Run air
air
