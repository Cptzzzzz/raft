#!/bin/bash

git pull

docker rmi raft

docker build -t raft .