#!/bin/bash
docker rmi -f $(docker images -q --filter "dangling=true")
