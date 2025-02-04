# utils

The utils folder contains small reusable models like:

* `decimal.proto` - A clean way to use decimals in the protos. Compatible with the `shopspring/decimal` package.
* `denom.proto` - Denom in a structure shape, capable of annotating the denoms in the coreum chain
* `metadata.proto` - Many proto message require the same metadata (network, created at, updated at). This is a way to avoid repeating the same fields in every message.

Functionality to manipulate these types is mostly generated, however a function to for example parse the denom into the denom type, would live in this folder.