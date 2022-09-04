DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

#1 : arm64v8 | arm32v7 | buster

num_arg=$#
if [ $num_arg -ne 1 ]; then 
    echo "Incorrect call: ./build.sh <target>"
    echo "example: ./build.sh arm64v8"
    exit 1
fi

target=$1
server_url=ghcr.io/tapvanvn/repository

pushd "$DIR"

dockerfile="docker/$target.dockerfile"

docker_image="ws-pubsub"

if test -f "$dockerfile"; then

    tag="$(<./version.txt)-$target"

    docker build -t $server_url/$docker_image:$tag -f $dockerfile ./

    docker push $server_url/$docker_image:$tag
else
    echo "The target $target is not supported"
    exit 1
fi

popd