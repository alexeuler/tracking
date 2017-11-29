#!/usr/bin/env bash

ab -c 10 -n 1000 -p postdata.json -H "X-Secret:test" -H "Origin:localhost" -T application/json http://staging.up-finder.com/events