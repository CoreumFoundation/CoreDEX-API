package main

import (
	"context"
	"math/rand"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	appConfig := AppConfig{
		GRPCHost: "full-node.devnet-1.coreum.dev:9090",
		// The token issuer used to issue the assets
		Issuer: AccountWallet{
			Address:  "devcore19p7572k4pj00szx36ehpnhs8z2gqls8ky3ne43",
			Mnemonic: "inmate connect object bid before sting talent interest forget tourist crystal girl estate banner cool crunch scatter industry sick motion hawk fossil seek slam",
		},
		// 2 accounts which are buying/selling the assets
		AccountsWallet: []AccountWallet{
			{
				Address:  "devcore1fksu90amj2qgydf43dm2qf6m2dl4szjtx6j5q8",
				Mnemonic: "carbon found inhale bitter sunny attack apple old hobby cave double dream priority north transfer visual select festival sunset fruit city increase empty rate",
			},
			{   
				Address:  "devcore1dj9yphkprdsuk6s4mgnfhnq5c39zf499nknkna",
				Mnemonic: "offer crop front arena tell because multiply glide cable claw goat sunset make tail bless race sound basket father pet across step wild occur",
			},
		},
		AssetFTDefaultDenomsCount: 2,
	}

	application, err := NewApp(ctx, appConfig)
	if err != nil {
		panic(err)
	}

	t := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case now := <-t.C:
			rootRnd := rand.New(rand.NewSource(now.Unix()))
			application.CreateOrder(ctx, rootRnd, application.GetAccounts())
		}
	}
}
