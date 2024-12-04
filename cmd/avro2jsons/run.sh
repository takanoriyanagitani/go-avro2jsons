#!/bin/sh

#export ENV_SCHEMA_FILENAME=./sample.d/sample.avsc
#cat sample.d/sample.jsonl | json2avrows > ./sample.d/sample.avro

cat sample.d/sample.avro |
	./avro2jsons
