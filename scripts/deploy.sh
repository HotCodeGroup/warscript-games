#!/bin/bash

set -e

echo -e "# Building docker.\n"
docker build -t warscript-games .
docker tag warscript-games $DOCKER_USER/warscript-games
docker push $DOCKER_USER/warscript-games

echo -e "# Starting docker.\n"

chmod 600 ./2019_1_HotCode_id_rsa.pem
ssh-keyscan -H 89.208.198.192 >> ~/.ssh/known_hosts
ssh -i ./2019_1_HotCode_id_rsa.pem ubuntu@89.208.198.192 docker pull $DOCKER_USER/warscript-games
for (( c=1; c<=$CONTAINERS_COUNT; c++ ))
do
    ssh -i ./2019_1_HotCode_id_rsa.pem ubuntu@89.208.198.192 docker stop warscript-games.$c
    if ssh -i ./2019_1_HotCode_id_rsa.pem ubuntu@89.208.198.192 test $? -eq 0
    then
        ssh -i ./2019_1_HotCode_id_rsa.pem ubuntu@89.208.198.192 docker rm warscript-games.$c || true
    fi
    ssh -i ./2019_1_HotCode_id_rsa.pem ubuntu@89.208.198.192 docker run -e CONSUL_ADDR=$CONSUL_ADDR \
                                                                    -e VAULT_ADDR=$VAULT_ADDR \
                                                                    -e VAULT_TOKEN=$VAULT_TOKEN \
                                                                    --name=warscript-games.$c \
                                                                    -d --net=host $DOCKER_USER/warscript-games
done