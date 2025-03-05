#!/bin/bash
set -e 

echo "Running SODA"
echo ""
echo ""

echo "  Migrating database:"

./migrate -verbose -database 'mysql://root:password@tcp(soda_db)/soda' -path db/migrations up & PID=$!
# Wait for migration to finish
echo "    waiting..."
wait $PID

echo " Files for debugging"
ls -lah

echo " Running soda executable"
echo ""
echo ""
./soda