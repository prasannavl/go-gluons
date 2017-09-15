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
    setup_bin_name
    BUILD_TARGET="${BUILD_TARGET:-${BIN_NAME}}"
    ASSETS_DIR="./www"
    BUILD_PACKAGE="${BUILD_TARGET/%.exe}.tar.gz"
    RUN_SCRIPT="./scripts/run.sh"
    DEPLOY_SERVER="labs.prasannavl.com"
    DEPLOY_TARBALL_DIR="apps.tarballs"
    DEPLOY_DIR="apps/${BUILD_TARGET}"
    if [[ "$GOOS" == "linux" ]]; then 
        BUILD_TARGET="${BUILD_TARGET/%.exe}"
    fi;
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
    tar -zvcf "$BUILD_PACKAGE" "$BUILD_TARGET" "$ASSETS_DIR" $RUN_SCRIPT
    echo "> pack: done"
}

deploy() {
    init
    if [ ! -f "$BUILD_PACKAGE" ]; then
        pack
    fi
    echo "> deploy: start"
    ssh $DEPLOY_SERVER mkdir -p "${DEPLOY_TARBALL_DIR}" "${DEPLOY_DIR}"
    rsync -avu "$BUILD_PACKAGE" $DEPLOY_SERVER:"$DEPLOY_TARBALL_DIR" --progress
    ssh $DEPLOY_SERVER "
    rm -rf \"${DEPLOY_DIR}/*\" &&
    tar -xvzf \"${DEPLOY_TARBALL_DIR}/${BUILD_PACKAGE}\" -C \"${DEPLOY_DIR}\" &&
    chmod +x \"${DEPLOY_DIR}/${RUN_SCRIPT}\" &&
    \"${DEPLOY_DIR}/${RUN_SCRIPT}\"
    " 
    echo "> deploy: end"
}

setup_bin_name() {
    if [ -n "$BIN_NAME" ]; then
        return
    fi;
    local cname=$(basename $(pwd))
    local cbin_name="${cname}"
    if [[ $(uname -s) == MINGW* ]]; then
        cbin_name="${cbin_name}.exe"
    fi
    BIN_NAME="${cbin_name}"
}

if [ -z "$@" ]; then
    run
else 
    "$@"
fi;
popd > /dev/null