#!/usr/bin/env bash

mpcli rest-server --chain-id mpchain --trust-node --laddr tcp://0.0.0.0:1317 > /dev/null &

mpd start