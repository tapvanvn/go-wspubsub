DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"



pushd "$DIR/../../"
source "config/env.sh"
go build
./go-wspubsub >&2
popd
