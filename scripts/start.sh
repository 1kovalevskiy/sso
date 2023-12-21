#!/bin/sh

/migrator --config-path="./config/config.yml" --migrations-path=./migrations
/app --config-path="./config/config.yml"