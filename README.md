<img src="https://bakinbacon.io/img/head_logo.png" width="400px">

---

## Running BakinBacon

_BakinBacon defaults to Granadanet, the current mainnet testing network. Use `-network mainnet` to switch._

1. Download the latest binary for your OS from [bakinbacon/releases](https://github.com/bakingbacon/bakinbacon/releases)
1. Execute: `./bakinbacon [-debug] [-trace] [-webuiaddr 127.0.0.1] [-webuiport 8082] [-network mainnet|granadanet]`
1. Open http://127.0.0.1:8082/ in your browser

## Building BakinBacon

If you want to contribute to BakinBacon with pull-requests, you'll need a proper environment set up to build and test BakinBacon.

### Dependencies

* Install go-1.16+
* Install nodejs-14.15 (npm 6.14)

### Build Steps

1. Clone the repo
1. `cd webserver; npm install && npm run build` (Build the webserver UI)
1. `go build` (This will download any go-lang dependencies and bundle the UI)
1. `./bakinbacon [-debug] [-trace] [-webuiaddr 127.0.0.1] [-webuiport 8082] [-network mainnet|granadanet]` (Be sure to go back one directory)
1. Open http://127.0.0.1:8082/ in your browser
