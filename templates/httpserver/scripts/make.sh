#!/usr/bin/env bash
## Author: Prasanna V. Loganathar

run() {
    init
    clean
    build
    test
}

init() {
    if [ -z "$__INIT" ]; then 
        __INIT="1"
    else
        return 0
    fi
    trap "echo '> script: incomplete termination requested'" TERM   
    set -e
    # script dir
    local dir=$( dirname "${BASH_SOURCE[0]}" )
    # Go to the parent directory of this script's dir
    local pdir="$( cd "${dir}/../" && pwd )"
    pushd . > /dev/null
    cd ${pdir}
    setup_vars
}

setup_vars() {
    PKG_BASE_NAME="$(basename $(pwd))"
    BUILD_TARGET="${BUILD_TARGET:-${PKG_BASE_NAME}}"
    ASSETS_DIR="./assets"
    
    BUILD_PACKAGE="${BUILD_TARGET}.tar.gz"
    RUN_SCRIPT="./scripts/run.sh"

    DEPLOY_SERVER="labs.prasannavl.com"
    DEPLOY_TARBALL_DIR="apps.tarballs"
    DEPLOY_DIR="apps/${BUILD_TARGET}"
}

clean() {
    init
    echo "> clean: start"    
    go clean
    rm "$BUILD_PACKAGE" || true
    echo "> clean: done"
}

build() {
    init
    echo "> build: start"
    go build -o "${BUILD_TARGET}"
    echo "> build: done"
}

test() {
    init   
    echo "> test: start"
    go test
    echo "> test: done"
}

pack() {
    init
    if [ ! -f "$BUILD_TARGET" ]; then
        run
    fi
    echo "> pack: start"    
    tar -zvcf "$BUILD_PACKAGE" "$BUILD_TARGET" "$ASSETS_DIR" "$RUN_SCRIPT"
    echo "> pack: done"
}

deploy() {
    init
    if [ ! -f "$BUILD_PACKAGE" ]; then
        pack
    fi
    echo "> deploy: start"

    local pre_script="
    mkdir -p \"${DEPLOY_TARBALL_DIR}\" \"${DEPLOY_DIR}\"
    "
    
    local post_script="
    rm -rf \"${DEPLOY_DIR}/{${BUILD_TARGET},${ASSETS_DIR}}\" &&
    tar -xvzf \"${DEPLOY_TARBALL_DIR}/${BUILD_PACKAGE}\" -C \"${DEPLOY_DIR}\" &&
    chmod +x \"${DEPLOY_DIR}/${RUN_SCRIPT}\" &&
    \"${DEPLOY_DIR}/${RUN_SCRIPT}\"
    "

    ssh $DEPLOY_SERVER "${pre_script}"
    rsync -avu "$BUILD_PACKAGE" $DEPLOY_SERVER:"$DEPLOY_TARBALL_DIR" --progress    
    ssh $DEPLOY_SERVER "${post_script}"
    echo "> deploy: end"
}

cleandeploy() {
    run
    deploy
}

if [ -z "$@" ]; then
    run
else 
    "$@"
fi;
popd > /dev/null