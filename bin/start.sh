#!/bin/bash

cd apps/store
./bin/start.sh &
# Give some time for the store to be ready (it can be creating (new) indexes, so start can take a while depending on the database size)
sleep 30
cd ../data-aggregator
./bin/start.sh &
cd ../api-server
./bin/start.sh &
cd ../frontend
./bin/start.sh &

jobs
echo "Servers started. Navigate to http://localhost:3000 to access the web interface."