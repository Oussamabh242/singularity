#!/bin/bash

# Define the number of clients
NUM_CLIENTS=100

# Loop through the number of clients
for i in $(seq 1 $NUM_CLIENTS); do
    # Start each client in the background and pass the client number
    python3 test.py bugs "$i" &
done

# Wait for all background processes to finish (optional)
wait
echo "All clients have been started."
