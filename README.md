# EScompanion

A companion wrapper for elasticsearch which can install Elasticsearch and Kibana plugins, compatible with all major elasticsearch versions starting v1.6 to v2.3.

## How to use

The EScompanion depends on `.env` file to get the credentials if you want
to install plugins on a remote system. If no `.env` file is provided then
EScompanion will revert to local system.

```bash
go build

# will install on the local system
esplugin -version 2.3 appbaseio/dejaVu mobz/elasticsearch-head
```

## How to test

The tests assume that the `plugin` binary is present in the location
 `/usr/share/elasticsearch/bin` if thats not true for your server then
 just change the command from the `testCommandProvider` function present in
 the test file.

Create a .env file which contains
```
USER=Test  # user on the server
SERVER=Elasticsearch.local
PORT=22
PASSWORD=neverTrustMe
    or
KEY="./.ssh/id_rsa"
```

Then `go test`
