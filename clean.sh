#!/bin/bash

set -e

rm -rf ./data/*.json
rm -rf ./data/out
rm -rf ./data/err
rm -rf ./data/renamed
rm -rf ./*.log

# rm -rf ./rename ./preproc