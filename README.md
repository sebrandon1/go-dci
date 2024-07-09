# go-dci

## Overview

A Golang based wrapper around the Red Hat Distributed CI API:

https://doc.distributed-ci.io/dci-control-server/docs/API/

## CLI Usage

Build the binary:

`make build`

Set the configuration:

```
$ ./go-dci config set -h
Set a key value pair to the configuration

Usage:
  dci config set [flags]

Flags:
  -a, --accesskey string   The access key to set in the configuration.
  -h, --help               help for set
  -s, --secretkey string   The secret key to set in the configuration.
```

The configuration accepts both `--accesskey` and `--secretkey` fields.  These refer to the DCI RemoteCI credentials.

Query the `jobs` from the API:

```
$ ./go-dci jobs -d 30
Getting all jobs from DCI that are 30 days old
Job ID: 78ed13e1-841f-4c04-a1c6-8df9028c67cd  -  TNF Version: tnf-v5.1.3 (Days Since: 11.249845)
Job ID: eb491abd-ec8b-42cc-aa8b-98adf741b236  -  TNF Version: tnf-v5.1.3 (Days Since: 11.305408)
Job ID: 6d424176-475a-4ec9-9f6c-7d43b9d3315b  -  TNF Version: tnf-v5.1.1 (Days Since: 13.123290)
Job ID: 56074acb-d21d-46eb-92f3-93c614b4ab67  -  TNF Version: tnf-v5.1.1 (Days Since: 13.198050)
Job ID: 27db50ce-dace-4709-ba45-9ae375918420  -  TNF Version: tnf-v5.1.1 (Days Since: 14.225086)
Job ID: f5a29537-a2b6-4312-a54d-e5da7e60114a  -  TNF Version: tnf-v5.1.1 (Days Since: 18.182748)
Job ID: 821ad00a-2cf5-4875-b308-4e37cbe8b7ab  -  TNF Version: tnf-v5.1.1 (Days Since: 28.121251)
```

Note: These jobs only pertain to those running the [test-network-function/cnf-certification-test](https://github.com/test-network-function/cnf-certification-test) suite.
