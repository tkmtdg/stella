FROM golang:1.11.3

WORKDIR /go/src/github.com/tkmtdg/stella

RUN apt-get update && apt-get install -y \
  zip \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*
RUN go get -u github.com/golang/dep/cmd/dep

CMD /bin/bash package.sh
