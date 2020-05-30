#!/bin/bash

if [[ $# -eq 0 ]] ; then
    uuid=$(uuid)
    aws sqs send-message --queue-url https://sqs.eu-west-1.amazonaws.com/041936244769/csqsv2.fifo --message-body "test-$uuid" --message-group-id "irrlab"
else
    aws sqs send-message --queue-url https://sqs.eu-west-1.amazonaws.com/041936244769/csqsv2.fifo --message-body "$@" --message-group-id "irrlab"
fi
