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

	appConfig := app.AppConfig{
		GRPCHost: "full-node.devnet-1.coreum.dev:9090",
		Issuer: app.AccountWallet{
			Address:  "devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs",
			Mnemonic: "final warrior tell admit apology road unlock gadget east airport clever roast whale ability lecture audit slot betray rapid legal crumble receive distance bind",
		},
		AccountsWallet: []app.AccountWallet{
			// {
			// 	Address:  "devcore180amp4nx5c4t00kk9z3f66lkcspz3zcadzvtfu",
			// 	Mnemonic: "fun like brave brand holiday air record aim angle permit reject sign wealth stumble quiz crew hunt trophy banana ritual hover minimum worth vendor",
			// },
			{
				Address:  "devcore16l8kr57njhzzn2gnp52vy96fnhsamnj0t2mw9t",
				Mnemonic: "toddler finger sell grain unhappy until pool trust salute suspect solar mail banana essence crash adjust enlist hockey imitate situate festival gold calm pioneer",
			},
			{
				Address:  "devcore16f300zup3y4y8f3mht5sfkhljuxtftw6ypa7sd",
				Mnemonic: "food coconut twelve brief valley wild nurse dog essence salute govern venture pill notice skin reduce blade donate mutual north excite much radio second",
			},
			// {
			// 	Address:  "devcore1fpdgztw4aepgy8vezs9hx27yqua4fpewygdspc",
			// 	Mnemonic: "allow coral siren couch melody wool wall approve clarify emotion swear ensure ketchup climb sunset churn theme bean blame power destroy angle certain better",
			// },
			// {
			// 	Address:  "devcore1878pk82zlndhldglx26r606qcd886562mad59y",
			// 	Mnemonic: "silk loop drastic novel taste project mind dragon shock outside stove patrol immense car collect winter melody pizza all deputy kid during style ribbon",
			// },
			// {
			// 	Address:  "devcore1vv97hnfxtar7a87szstpdlzexwu009cylfdsry",
			// 	Mnemonic: "illness drift profit gravity asset call purse scan chef name shine odor badge curve stumble business affair come capable correct core indoor cargo dinner",
			// },
			// {
			// 	Address:  "devcore1yegtec4a0zd9ksepr6hf8au95940p5hvy06064",
			// 	Mnemonic: "phrase conduct toe genius ramp farm fuel tomato spider seek accident acoustic tower injury token never rate hire drum notable garage term dream flower",
			// },
			// {
			// 	Address:  "devcore12lyzywqkpj6jlmeujel32jzt4qu65q53mc9j2p",
			// 	Mnemonic: "old whale start lecture perfect load rhythm camera negative hero skirt scale they thank dwarf okay lawsuit elevator dry champion top bright myth theme",
			// },
			// {
			// 	Address:  "devcore19dk2skhjeurar3x0xrah4xeuesrntz8n6cjn86",
			// 	Mnemonic: "victory civil corn two axis luxury renew winter dry solve brief enter maze door slim fiction card trumpet slam attend flight orphan tank organ",
			// },
			// {
			// 	Address:  "devcore1mcgamu8yu9uan4qafq7pqa4e3heapsj6gqmnky",
			// 	Mnemonic: "jazz black mix busy bulk orbit blanket blur host sell wasp twin tobacco stomach secret patch present install february ecology ripple maximum chronic ridge",
			// },
		},
		AssetFTDefaultDenomsCount: 12,
	}

	application, err := app.NewApp(ctx, appConfig)
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
			ordersCount := 10
			accounts := application.GetAccounts()
			wg := sync.WaitGroup{}
			wg.Add(len(accounts))
			for _, account := range accounts {
				err = application.CreateOrder(ctx, rootRnd, account, ordersCount)
				wg.Done()
				if err != nil {
					panic(err)
				}
			}
			wg.Wait()
		}
	}
}
