#!/bin/bash -xv

set -e

#_DEBUG="on"
_DEBUG="off"

function DEBUG() {
    [ "$_DEBUG" == "on" ] && $@
}

export GIT_SSL_NO_VERIFY=1
export https_proxy=http://git.iizone.com.cn:32222 http_proxy=http://git.iizone.com.cn:32222 all_proxy=socks5://git.iizone.com.cn:32222

if [ ! -d "/opt/wrk" ]; then
    cd /opt/wrk
    git clone https://github.com/wg/wrk.git
    cd wrk
    make
    if [ $? -ne 0 ]; then
        echo "make wrk failed"
        exit 1
    fi
    if [ ! -f "/usr/local/bin/wrk" ]; then
        ln -s /opt/wrk/wrk/wrk /usr/local/bin/wrk
    fi
else
    cd /opt/wrk/wrk
    git pull
    make
    if [ $? -ne 0 ]; then
        echo "make wrk failed"
        exit 1
    fi
    if [ ! -f "/usr/local/bin/wrk" ]; then
        ln -s /opt/wrk/wrk/wrk /usr/local/bin/wrk
    fi
fi

echo "build wrk success"

