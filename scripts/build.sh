SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
SRC_DIR="${SCRIPT_DIR}/../"

go build -o "$SRC_DIR"/build/client "$SRC_DIR"/cmd/