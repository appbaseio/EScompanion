# ESplugin

A plugin wrapper for elasticsearch which works for all elasticsearch versions.

## How to use

The ESplugin depends on `.env` file to get the credentials if you want
to install plugins on a remote system. If no `.env` file is provided then
ESplugin will revert to local system.

```bash

go build

#will install on the local system
esplugin -version 2.3 appbaseio/DejaVu mobz/elasticsearch-head


```

## How to test

Create a .env file which contains
```
User=Test  # user on the server
Server=http://Elasticsearch.local
Port=22
Password=neverTrustMe
    or
Key="./.ssh/id_rsa"
```

Then `go test` 
