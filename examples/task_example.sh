#!/bin/bash

# This is an example task that simulates some work
# It will:
# 1. Print a message
# 2. Sleep for a random duration (1-5 seconds)
# 3. Generate some random data
# 4. Exit with success or failure randomly

echo "Starting example task at $(date)"
echo "Task parameters: $@"

# Sleep for random duration
SLEEP_TIME=$((RANDOM % 5 + 1))
echo "Working for $SLEEP_TIME seconds..."
sleep $SLEEP_TIME

# Generate some random data
echo "Generating random data..."
for i in {1..5}; do
    echo "Data point $i: $RANDOM"
done

# Randomly succeed or fail
if [ $((RANDOM % 10)) -lt 8 ]; then
    echo "Task completed successfully at $(date)"
    exit 0
else
    echo "Task failed at $(date)"
    exit 1
fi 