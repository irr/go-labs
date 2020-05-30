#!/bin/bash

# parallel -n0 -j100 ./msg.sh ::: {0..100}

aws sqs get-queue-attributes --queue-url https://sqs.eu-west-1.amazonaws.com/041936244769/csqsv2.fifo --attribute-names All
