#!/bin/bash

# Enter development continer

if [[ $OSTYPE =~ ^msys || $OSTYPE =~ ^win ]]; then
    # NOTE: Not tested on windows
    # Use 'winpty' to start an interactive shell in the container
    winpty docker exec -it seegoals_controller-1 sh
else
    docker exec -it seegoals-controller-1 sh
fi
