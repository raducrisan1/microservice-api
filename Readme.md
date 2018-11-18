# The API Microservice

## Overview

The API microservice fetches its data from Reports. It uses gRPC for this purpose and the message syntax is defined in the file `tradesuggest.proto`. To generate the proxy and stub code around this, use the protoc tool like below:

```bash
protoc -I=tradesuggest --go_out=plugins=grpc:tradesuggest tradesuggest/tradesuggest.proto
```

Note: you need to have protoc-gen-go plugin installed together with protoc tool. For further information, please consult: [GRPC Go Quick Start](https://grpc.io/docs/quickstart/go.html)

## Creating the Docker container

To create a docker container, simply run 
