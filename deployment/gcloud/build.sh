# we assume that 

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

pushd "$DIR/../../"

tag=$(<./version.txt)

server_url=tapvanvn

docker build -t $server_url/ws_pubsub:$tag -f docker/gcloud.dockerfile ./

docker push $server_url/ws_pubsub:$tag

popd