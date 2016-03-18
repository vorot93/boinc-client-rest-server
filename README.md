# RESTful API for BOINC client

Access your BOINC client via HTTP. Implements the unofficial [BOINC RPC spec](https://github.com/vorot93/boinc-undocumented).

## Installation and Running
### Installation

```
$ make get-deps
$ make build
```

### Running

```
$ ./build/boinc-client-rest-shim
```

## Using the API

```
$ curl -d 'json={"host": <BOINC client RPC address>, "pwd": "<BOINC client RPC key>"}' localhost:15678/0.1/<entrypoint>
```

## Available entrypoints

- `/tryauth` - test authentication.
- `/projects` - view all projects known to BOINC client.
- `/messages` - show BOINC client logs.
- `/acctmgr` - poll account manager manipulation status.
- `/acctmgr/info` - show information about current account manager.
- `/quit` - shut down the shim.
- 
