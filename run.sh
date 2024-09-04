#!/bin/bash	

#----------------------------
#TODO: is it: independen fo all config files, databases, login providers?
make build
#----------------------------

TARGET_DIR=./cmd/app/
OUTPUT_BINARY=console-api


# region: us-west-2
export AWS_ACCESS_KEY_ID=AKIA5JJVZHBG537GSSQ2
export AWS_SECRET_ACCESS_KEY=p+/NstwmLF7xvfzL2qvh/tnlo1mFsfAmaYKnK/F4
# export AWS_SESSION_TOKEN=your_session_token # if you are using temporary credentials

export LOGINPROVIDER=cognito
export DBTYPE=pg
export PGCASBINURI=postgres://postgres:okokokokd@localhost:5432/casbin
export PGWEATHERURI=postgres://postgres:okokokokd@localhost:5432/weather

cd ${TARGET_DIR}
./${OUTPUT_BINARY} 
