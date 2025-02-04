# Rates

The rates package uses coingecko, the pools and the Dijkstra algorithm to resolve any currency to USDC (assuming there is a path).

## Start parameters

Required is the `BASE_COIN` which indicates the ultimate root (SARA) in case we can not resolve the currency to USDC using the Dijkstra algorithm and the present pools.
Also required is the `BASE_USDC` which is the USDC coin in the network, so the coin to which Dijkstra will resolve the currency to.

- `BASE_COIN` - Structure `{"BaseCoin":[{{"Network": "devnet","Coin": "usara-devcore1wkwy0xh89ksdgj9hr347dyd2dw7zesmtrue6kfzyml4vdtz6e5wsyjwwgp"}]}`
- `BASE_USDC` - Structure `{"BaseCoin":[{"Network": "mainnet","Coin": "uusdc-E1E3674A0E4E1EF9C69646F9AF8D9497173821826074622D831BAB73CCB99A2D"}]}`