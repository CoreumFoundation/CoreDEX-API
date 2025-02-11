# Stop docker compose and remove volumes 
# (just in case there is a new init script which else would be skipped in a next run)
docker compose down -v