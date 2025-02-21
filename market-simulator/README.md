# Market simulator

The market simulator is a 2-user trading app that provides continuous trading in 1 market. The app funds the accounts on a regular basis so that the simulation can continue indefinitely.

The app places orders on dex test denom0 and 1 between 80 and 100 in value, and an amount of 10 to 20. The orders are valid for 1 hour so that the order limit (100 orders max for a market) is avoided or managed. The app places 1 order per minute.

## Get up and running on localhost

The application is in golang and requires virtually nothing. The app uses the blockchain only, and is acting against devnet.

A typical start is:

```bash
go run .
```
