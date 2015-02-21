#!/bin/bash
LIMIT=${LIMIT:-5}
ARGS="-limit=$LIMIT"
if [ ! -z "$VERBOSE" ]; then
    ARGS="$ARGS -verbose"
fi
if [ ! -z "$PERSISTENT" ]; then
    ARGS="$ARGS -file=/var/resm.db"
fi
CMD="/src/bin/resm $ARGS"
#echo $CMD
$CMD
