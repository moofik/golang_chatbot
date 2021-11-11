#!/bin/bash

pg() {
  local base_dir=$(dirname ${DIR})
  local work_dir=$(pwd)

  if [[ ${work_dir} == $(pwd) ]]; then
    work_dir="${PROJECT_DIR}"
  fi

  ${SUDO_CMD} docker run \
    -it \
    --rm \
    -v ${DIR}:${PROJECT_DIR} \
    -w ${work_dir} \
    -e POSTGRES_PASSWORD=example \
    --network crypto \
    crypto/postgresql \
    "$@"
}