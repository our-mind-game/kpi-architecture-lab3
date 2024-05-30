#!/bin/bash

startx=100
starty=100
steps=5
interval=1

sleep $interval
curl -X POST http://localhost:17000 -d "reset
white
figure $startx $starty"

for (( i=1; i<=steps; i++ ))
do
    newx=$((startx + i))
    newy=$((starty + i))
    
    curl -X POST http://localhost:17000 -d "move $newx $newy
    update"
    
    sleep $interval
done

