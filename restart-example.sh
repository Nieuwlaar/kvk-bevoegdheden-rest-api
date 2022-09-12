#!/bin/sh
go build $PWD
killall kvk-extract-manual
CERTIFICATE_KVK="-----BEGIN CERTIFICATE-----
xxx
xxx
xxx
xxx
-----END CERTIFICATE-----" PRIVATE_KEY_KVK="-----BEGIN RSA PRIVATE KEY-----
xxx
xxx
xxx
-----END RSA PRIVATE KEY-----" ./kvk-extract-manual &