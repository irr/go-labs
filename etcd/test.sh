#!/bin/bash
http --follow --form POST http://127.0.0.1:4001/v1/keys/foo value=ivan
http http://127.0.0.1:4001/v1/keys/foo
