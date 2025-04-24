# Market simulator

The market simulator is a 2-user trading app that provides continuous trading in 1 market. The app funds the accounts on a regular basis so that the simulation can continue indefinitely.

The app places orders on dex test denom0 and 1 between 80 and 100 in value, and an amount of 10 to 20. The orders are valid for 1 hour so that the order limit (100 orders max for a market) is avoided or managed. The app places 1 order per minute.

## Get up and running on localhost

The application is in golang and requires virtually nothing. The app uses the blockchain only, and is acting against devnet.

A typical start is:

```bash
go run .
```

## Application start parameters

- `APP_CONFIG` - The application config

### APP_CONFIG

The app config is defined such that parties can run their own simulators. Since the simulator itself creates the tokens, the first address is the issuer of the token. Subsequent addresses are the accounts that will be funded and used for trading.

```json5
{
    "GRPCHost": "full-node.devnet-1.coreum.dev:9090",
    "Issuer": {
        "Address": "devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
        "Mnemonic": "inmate connect object bid before sting talent interest forget tourist crystal girl estate banner cool crunch scatter industry sick motion hawk fossil seek slam"
    },
    "AccountsWallet": [
        {
            "Address": "devcore1fksu90amj2qgydf43dm2qf6m2dl4szjtx6j5q8",
            "Mnemonic": "carbon found inhale bitter sunny attack apple old hobby cave double dream priority north transfer visual select festival sunset fruit city increase empty rate"
        },
        {
            "Address": "devcore1dj9yphkprdsuk6s4mgnfhnq5c39zf499nknkna",
            "Mnemonic": "offer crop front arena tell because multiply glide cable claw goat sunset make tail bless race sound basket father pet across step wild occur"
        }
    ],
    "AssetFTDefaultDenomsCount": 2
}
```