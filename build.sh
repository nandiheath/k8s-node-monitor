VERSION=0.1.1

docker build -t kube-node-monitor .

docker tag kube-node-monitor nandiheath/kube-node-monitor:$VERSION
docker push nandiheath/kube-node-monitor:$VERSION