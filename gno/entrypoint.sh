#!/bin/sh
set -e

./build/gnoland config init

./build/gnoland config set rpc.laddr tcp://0.0.0.0:${NODE_PORT}

exec ./build/gnoland start --lazy