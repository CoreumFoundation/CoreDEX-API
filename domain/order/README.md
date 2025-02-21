# Order

The Order proto provides all the functionality required to interact with the order store.
The Order proto also defines the domain objects for communication between the services and use events (e.g., messaging).

The Order is modeled after the order as defined in the Coreum DEX.

## Building the required files

Once the proto file is updated, the following files need to be generated:

* go

### go

There is a file, proto.sh, which can be used to generate the go files.

```sh
./bin/build.sh
```
