# CoreDEX API: The API for the Coreum DEX, A business-ready exchange

## About the Coreum DEX

The Coreum DEX is a permissionless, fee-free exchange with an ***organic order book*** designed for RWAs. Build trading interfaces using the CoreDEX APi pre-built extensions for compliance and increased security.

| Quick Setup                                                     | Custom Trading Logic                                                                                      | Developer Ready API                                                                                                         |
|-----------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------|
| Get started in under 15 minutes with minimal technical overhead | Configure trading hours, compliance rules, and asset restrictions with built-in smart contract extensions | Use any language to interact with the API, providing trade history, order history, OHLC data and some convenience endpoints |


## About the CoreDEX API

The CoreDEX is an API on top of the Coreum blockchain integrated DEX. The API is in RESTFUL style so that it is easier to use than using the blockchain directly. By using this restful API, the developer only needs to integrate signing of transactions with their preferred wallet system.

The project provides a data aggregator (blockchain scanner), an API server, store, and a UI (for demo purposes).

Once started, the UI is accessible at `http://localhost:3000`.

The API provides (folder `/apps/api-server`)

* Trade history
* Order tracking
* OHLC (Open/High/Low/Close) data
* Tickers
* Order entry and execution

Furthermore, the CoreDEX API provides data aggregation for the trades and orders, thus generating the OHLC (folder `/apps/data-aggregator`).

The API server is horizontally scalable, thus providing the developer/operator with the ability to serve 1000s of parallel requests.

The data-aggregator is a single instance and is build such that it will, in case of a crash, restart where it left of before the crash, thus making certain the full history is present.

The store (folder `/apps/store`) is an implementation against a MySQL database, which in the end is the suspected bottleneck for any scaling of the system. The store is connected with the API server and the data-aggregator through grpc, and can thus be replaced with a different store implementation if scaling requires this.

## Installation and running of the CoreDEX API

There are 3 provided methods of running the Friendly DEX:

* Localhost with handy shell start scripts: Good for development and testing
* Docker-compose: Good for testing and development
* Kubernetes example deployment files: Handy for production

### Get up and running on localhost

Getting the CoreDEX running on your local (unix or linux) system takes just a few minutes.

To run the CoreDEX, you need to have the following installed:

* Go lang 1.23 or higher (https://golang.org/dl/)
* MySQL (https://dev.mysql.com/downloads/)

In MySQL create a database called `friendly_dex` and a user `testuser` with password `password`:

```sql
CREATE DATABASE friendly_dex;
CREATE USER 'testuser'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON friendly_dex.* TO 'testuser'@'localhost';
FLUSH PRIVILEGES;
```

There is a start script for all the components, excluding MySQL, in the `bin` directory of each component and in the main `/bin` directory.

Run

```bash
./bin/start.sh
```

### First start

On the first start, the process will start scanning from block 1 in the chain. So catching up is going to take a while.
If you, however, want to get to just current, you can set the state in the database to the current block number.

Steps to do this are:

1. Stop the data-aggregator
2. Run the following SQL command:

```sql
UPDATE State SET Content = '{"Height":6618678}' WHERE StateType=1;
```

After which the data-aggregator can be started again by running the start script.

### Inspecting running processes

Using the `ps` command you can see the running processes:

```bash
ps -ef | grep "go run"
```

You should see the following processes:

* `data-aggregator` - the data aggregator
* `api-server` - the API server
* `store` - the store

### Optional requirements

The configuration provided uses the devnet public node of the Coreum blockchain. While there is a public node for Coreum available, that node does not have a guaranteed performance and the installation of a local node is recommended.
See the [Coreum](https://docs.coreum.dev/docs/next/nodes-and-validators/run-full-node) for more information on how to install a local node.

### Running the Friendly DEX with Docker-compose

There is a docker compose file, `docker-compose.yml`, in the root of the project. This file can be used to start the CoreDEX API in a docker environment.

***Important***:
To use the docker compose environment, edit the folder `mysql-init/init-databases.sql` to set the block height to start from.
Once that has been done, the start script can be run:

```bash
./bin/start-docker-compose.sh
```

The script, `./bin/stop-docker-compose.sh`, can be used to stop the docker-compose environment. This script will remove the volumes (e.g., reset the database), so depending on your use case you might want to just shutdown with `docker compose down`.

### Running the Friendly DEX on VMs, Kubernetes or other cloud services

The architecture of the CoreDEX API such that it is horizontally scalable. Setup deployments such that:

* API Server and store are in stateless deployments, scale 1 to n
* Store can be a stateless deployment, scale 1 to n
* Data aggregator can be in stateless or stateful, scale 1 max

CPU and memory expectations:

* API Server: 1 CPU, 128MB memory
* Store: 1 CPU, 128MB memory
* Data aggregator: 1 CPU, 128MB memory
* UI: 1 CPU, 256MB memory

A set of single instances is expected to be able to handle 1000s of parallel requests.

To support deployment on cloud services, there are 4 docker files in the root of the project. These can be used to build the docker images for the components.

```bash
docker build -t coreumfoundation/api-server:latest -f Dockerfile.api-server .
docker build -t coreumfoundation/data-aggregator:latest -f Dockerfile.data-aggregator .
docker build -t coreumfoundation/store:latest -f Dockerfile.store .
docker build -t coreumfoundation/frontend:latest -f Dockerfile.frontend .
```

There are also sample deployment files for Kubernetes in the `apps/kubernetes` directory.

## Market simulator

For purposes of demonstration, there is a market simulator in the `apps/market-simulator` directory. This simulator can be used to simulate trading on the CoreDEX API. The checked in mnemonics are there for demonstration purposes only and should not be used in a production environment.

To build the market simulator, run:

```bash
docker build -t coreumfoundation/market-simulator:latest -f Dockerfile.market-simulator .
```

## Application documentation

See the `README.md` files in the associated application folders.
Do note the UI `README.md` where there is information on some javascript packages/files that take care of UI data management.

## Common issues/FAQ

### Blockheight not updating

There is an error in the data-aggregator log stating:

```
2025-05-20 11:35:10 {"level":"error","severity":"ERROR","time":"2025-05-20T15:35:10Z","message":"rpc error for block 1: error in json rpc client, with http response metadata: (Status: 200 OK, Protocol HTTP/1.1). RPC error -32603 - Internal error: height 1 is not available, lowest height is 36782260"}
```

This is due to the State on being set before the data-aggregator is started. Set the height manually for the given network by logging in into the database, verify the existence of the `State` table, and run the following SQL command:

```sql
INSERT INTO State (StateType, Content, MetaData) VALUES (1, '{"Height":36782260}', '{"Network": 3, "UpdatedAt": {"seconds": 1738799304, "nanos": 164479000}}');
```

Or if the record is present, update the height:

```sql
UPDATE State SET Content = '{"Height":36782260}' WHERE StateType=1;
```
