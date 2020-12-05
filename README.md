## Dolphin-Send

This program continually tails a log file and sends the desired messages to the Dolphin bot (or technically where ever you want) via HTTP as JSON messages.

## [![Report](https://goreportcard.com/badge/github.com/EbonJaeger/dolphin-send)](https://goreportcard.com/report/github.com/EbonJaeger/dolphin-send) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Building

Dolphin has a Makefile to make building and installing easier.

To build the project, run `make`. To check the project and run tests, run `make check`.

## Installation

Coming soon!

## Usage

```
./dolphin-send [OPTIONS]
```

Options:

```
-a, --address  - The address of the listening server to send to
    --debug    - Print additional debug lines to stdout
-h, --help     - Print the help message
-k, --keywords - Comma-separated list of words to use as additional death keywords to look for
-l, --log      - The location of the log file to read
-p, --port     - The port number of the listening server
```

## License

Copyright © 2020 Evan Maddock (EbonJaeger)  
Makefile adapted from [usysconf](https://github.com/getsolus/usysconf), which is a Solus Project

Dolphin-Send is available under the terms of the Apache-2.0 license
