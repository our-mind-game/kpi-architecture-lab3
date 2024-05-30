#!/bin/bash

curl -X POST http://localhost:17000 -d "reset
white
bgrect 25 25 775 775
figure 500 500
green
figure 250 250
update"
