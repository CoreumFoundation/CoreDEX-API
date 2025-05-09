#!/bin/bash

export LOG_LEVEL=info
export APP_CONFIG='{
    "Network": "devnet-",
    "Fund": "https://api.devnet-1.coreum.dev/api/faucet/v1/fund",
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
}'

go run .