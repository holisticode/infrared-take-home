# infrared-take-home

## Overview
This repository contains a submission to the take home task as specified at
https://gist.github.com/red4626/2b7dbe91d351e0a5fc25d89704fd38a7

## Goal
Implement an implementation for merkle proof generation of Randao data 

## Implementation 
The presented solution uses JSON files to represent beacon state and block root information,
which is parsed and then processed as MerkleTrees, of which then proofs are generated.

The test files are included in the repo.

To execute parsing and proof generation, run `go run main.go`
This will generate a proof for a hardcoded default randao index of 7. 
It is possible to pass it a custom value and the generate command:
`go run main.go generate 5` (just passing `go run main.go generate` will also use the default).

There is also a verification function:
`go run main.go verify <leaf> <proof> <root> <target>`. There is only limited argument sanity checks.
The arguments format expected:
    * leaf: hex-encoded hash
    * proof: comma-separated string list of hex-encoded hashes 
    * root: hex-encoded hash 
    * target: integer

To execute the tests, which also verify the proofs, run `go test -v ./pkg/...`

## Process 
My initial approach was to launch `geth` with a consensus client (`lodestar`),
and then query the consensus client to gather state and block information
for processing.

Unfortunately my attempts at this failed, I wasn't able to find working SSZ parsers in go
which were working directly on reading the API.

Finally I switched to use JSON. However, the node I was using was in some parts not returning
spec-compliant data (e.g. `previous_epoch_attestations` were called `previous_epoch_participations`,
and this fields were just lists of strings, not complex containers as in the spec).

So I made the decision to use static JSON files for implementation and testing.

## Limitations
This implementation would probably fail in a production environment.

Using JSON, and with no spec compliant `go` data representation of the data types 
in the specification, I implemented my own JSON conversion, which is using reflection.
To avoid too deep of a rabbit hole, I am converting `string` values from JSON fields to just their
`[]byte` representation, which probably doesn't comply with the spec, for example, when 
values are `uint64` or other numerical types.

A production ready implementation would have to address this limitation.

## Planned but not realized elements
I wanted to provide a full execution and testing environment by providing local
node abstractions or local runners.

Ideally I wanted also a local beacon chain development node.
Due to lack of knowledge of the availability of such plus lack of time,
I had to drop this idea.

I know I could have run a hardhat local execution client. 

The idea was to query the local beacon chain node, create proofs from that data,
and then verify these proofs inside the local execution client in a smart contract. 
