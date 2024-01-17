SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
SRC_DIR="${SCRIPT_DIR}/../"

cd $SRC_DIR && COMPOSE_PROJECT_NAME=seegoals docker compose up -d