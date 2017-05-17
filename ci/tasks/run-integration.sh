#!/bin/bash

set -e -x

[[ -s "/home/emc/.gvm/scripts/gvm" ]] >/dev/null 2>/dev/null
source "/home/emc/.gvm/scripts/gvm" >/dev/null 2>/dev/null

PROJECT_PATH=$GOPATH/src/github.com/RackHD
COMPOSE_PATH=ci/integration/docker-compose.yml

cleanUp()
{
  # Don't exit on error here. All commands in this cleanUp must run,
  #   even if some of them fail
  set +e

  # Delete all containers
  docker rm -f $(docker ps -a -q)

  # Delete all images
  docker rmi -f $(docker images -q)

  # Delete any dangling volumes
  docker volume rm $(docker volume ls -qf dangling=true)

  # Clean up all cloned repos
  cd $GOPATH
  rm -rf $GOPATH/src
}

trap cleanUp EXIT

pushd $PROJECT_PATH/voyager-secret-service
  echo "Testing Secret Service"

  docker-compose -f ${COMPOSE_PATH} create
  docker-compose -f ${COMPOSE_PATH} start

  make deps
  make integration-test

  docker-compose -f ${COMPOSE_PATH} kill
  # Delete all containers
  docker rm -f $(docker ps -a -q)

  echo "Secret Service PASS\n\n"
popd

exit
