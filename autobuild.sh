#!/bin/sh

# go insists on absolute path.
export GOBIN=`pwd`/dist
export DISTDIR=`pwd`/dist
export DEVOS=linux

# space seperated packages
PACKAGES=`cd cmd/cacheservice && ls -1 |grep -v vendor`

buildall() {
    echo building ${GOOS}/$GOARCH
    GOBIN=${DISTDIR}/${GOOS}/${GOARCH}
    mkdir -p $GOBIN
    for pkg in ${PACKAGES}; do
	MYSRC=cmd/cacheservice/${pkg}
	( cd ${MYSRC} && make all ) || exit 10
    done
}

if [ -d dist ]; then
    rm -rf dist
fi
mkdir dist

# we only build for amd64 atm
export GOARCH=amd64

# this allows local builds on -dev machines
# to quickly build only a single arch
# intent is for devs to set DEVOS=[localos] permanently
# on their machine and
# the autobuild.sh will do 'The Right Thing'
if [ ! -z "${DEVOS}" ]; then
    GOOS=${DEVOS}
    buildall
    exit 0
fi

#========= build linux
export GOOS=linux ; buildall

#========= build mac
export GOOS=darwin ; buildall

#========= build windows
export GOOS=windows ; buildall

wait
echo $?
