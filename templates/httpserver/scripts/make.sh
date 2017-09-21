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
    run
    pack
    deploy-package "$BUILD_PACKAGE"
}

deploy-package() {
    echo "> deploy: start"
    local package="$1"
    if [ ! -f "$package" ]; then
        echo "package not found: \"$package\""
        return
    fi
    local pre_script="
    mkdir -p \"${DEPLOY_TARBALL_DIR}\" \"${DEPLOY_DIR}\"
    "
    local post_script="
    rm -rf \"${DEPLOY_DIR}/{${BUILD_TARGET},${ASSETS_DIR}}\" &&
    tar -xvzf \"${DEPLOY_TARBALL_DIR}/${package}\" -C \"${DEPLOY_DIR}\" &&
    chmod +x \"${DEPLOY_DIR}/${RUN_SCRIPT}\" &&
    \"${DEPLOY_DIR}/${RUN_SCRIPT}\"
    "

    ssh "$DEPLOY_SERVER" "${pre_script}"
    rsync -avu "$package" "$DEPLOY_SERVER":"$DEPLOY_TARBALL_DIR" --progress    
    ssh "$DEPLOY_SERVER" "${post_script}"
    echo "> deploy: end"
}

if [ -z "$@" ]; then
    run
else 
    "$@"
fi;
popd > /dev/null