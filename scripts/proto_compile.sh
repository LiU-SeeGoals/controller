SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PROTO_DIR="${SCRIPT_DIR}/../internal/proto"

protoc -I $PROTO_DIR --go_out="$PROTO_DIR"/ssl_vision "$PROTO_DIR"/*.proto