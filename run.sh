#!/bin/bash

# Determine the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Assuming the script is located at the root of the project
PROJECT_ROOT="${SCRIPT_DIR}"

# Define the paths to the build script and the main.go file
BUILD_SCRIPT="${PROJECT_ROOT}/internal/proto/build_go.sh"
MAIN_GO="${PROJECT_ROOT}/cmd/main.go"

echo "${BUILD_SCRIPT}"

# Run the build script
chmod +x ${BUILD_SCRIPT}
bash ${BUILD_SCRIPT}

if [ $? -eq 0 ]; then
	echo "Build successful. Running Go application..."
	go run "${MAIN_GO}"
else
	echo "Build failed. Exiting..."
fi
