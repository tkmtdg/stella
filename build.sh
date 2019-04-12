#!/bin/bash
docker build . -t stella-build
docker run -v `pwd`:/go/src/github.com/tkmtdg/stella stella-build
