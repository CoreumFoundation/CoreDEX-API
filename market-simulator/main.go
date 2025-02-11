package main

import (
	"context"
	"math/rand"
	"os/signal"
	"sync"
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
			Address:  "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
			Mnemonic: "final warrior tell admit apology road unlock gadget east airport clever roast whale ability lecture audit slot betray rapid legal crumble receive distance bind",
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

	t := time.NewTicker(time.Second)

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case now := <-t.C:
			rootRnd := rand.New(rand.NewSource(now.Unix()))
			accounts := application.GetAccounts()
			wg := sync.WaitGroup{}
			wg.Add(len(accounts))
			for _, account := range accounts {
				application.CreateOrder(ctx, rootRnd, account)
				wg.Done()
			}
			wg.Wait()
		}
	}
}
