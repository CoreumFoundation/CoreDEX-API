# CoreDEX API

The CoreDEX is an API on top of the Coreum blockchain integrated DEX. The API is in RESTFUL style so that it is easier to use than using the blockchain directly. By using this restful API the developer onbly needs to integrate signing of transactions with their preferred wallet system.

The project provides a data aggrator (block chain scanner), an API server, store, and an UI (for demo purposes).

Once started, the UI is accessible at `http://localhost:3000`.

The API provides (folder `/apps/api-server`)

* Trade history
* Order tracking
* OHLC (Open/Hight/Low/Close) data
* Tickers
* Order entry and execution

Furthermore the Friendly DEX provides data aggregation for the trades and orders, thus generating the OHLC (folder `/apps/data-aggregator`).

The API server is horizontally scalable thus providing the developer/operator with the ability to serve 1000s of parallel requests.

The data-aggregator is single instance and is build such that it will, in case of a crash, restart where it left of before the crash, thus making certain the full history is present.

The store (folder `/apps/store`) is an implementation against a MySQL database, which in the end is the suspected bottleneck for any scaling of the system. The store is connected with the API server and the data-aggregator through grpc, and can thus be replaced with a different store implementation if scaling would require this.

## Installation and running of the Friendly DEX

Getting the CoreDEX running on your local (unix or linux) system takes just a few minutes.

To run the CoreDEX you need to have the following installed:

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
If you however want to get to just current, you can set the state in the database to the current block number. 

Steps to do this are:

1. Stop the data-aggregator
2. Run the following SQL command:

```sql
UPDATE State SET content = '{"Height":6618678}' WHERE StateType=1;
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
See the [Coreum](https://docs.coreum.dev/docs/become-validator/run-full-node) for more information on how to install a local node.

## Running the Friendly DEX on VMs, Kubernetes or other cloud services

The architecture of the Friendly DEX such that it is horizontally scalable. Setup deployments such that:

* API Server and store are in stateless deployments, scale 1 to n
* Data aggregator can be in stateless or stateful, scale 1 max

CPU and memory expectations:

* API Server: 1 CPU, 128MB memory
* Store: 1 CPU, 128MB memory
* Data aggregator: 1 CPU, 128MB memory
* UI: 1 CPU, 256MB memory

A set of single instances is expected to be able to handle 1000s of parallel requests.

To support deployment on cloud services, there are 4 docker files in the root of the project. These can be used to build the docker images for the components.

```bash
docker build -t api-server -f Dockerfile.api-server .
docker build -t data-aggregator -f Dockerfile.data-aggregator .
docker build -t store -f Dockerfile.store .
docker build -t frontend -f Dockerfile.frontend .
```

## Application documentation

See the `README.md` files in the associated application folders. 
Do note the UI `README.md` where there is information on some javascript packages/files which take care of UI data management.
