package main

import (
	"context"
	"testing"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/CoreumFoundation/CoreDEX-API/apps/api-server/app/order"
	"github.com/CoreumFoundation/CoreDEX-API/domain/currency"
	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	"github.com/CoreumFoundation/CoreDEX-API/domain/metadata"
	orderservice "github.com/CoreumFoundation/CoreDEX-API/domain/order"
)

func TestOrders(t *testing.T) {
	ctx := context.Background()

	denom1, err := denom.NewDenom("dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs")
	if err != nil {
		t.Fatal(err)
	}
	denom2, err := denom.NewDenom("udevcore")
	if err != nil {
		t.Fatal(err)
	}

	orderService := orderservice.NewMockOrderServiceClient()
	currencyService := currency.NewMockCurrencyServiceClient()
	p1 := int32(6)
	denom1.Precision = &p1
	_, err = currencyService.Upsert(ctx, &currency.Currency{
		Denom: denom1,
		MetaData: &metadata.MetaData{
			Network:   metadata.Network_DEVNET,
			UpdatedAt: timestamppb.Now(),
			CreatedAt: timestamppb.Now(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	p2 := int32(0)
	denom2.Precision = &p2
	_, err = currencyService.Upsert(ctx, &currency.Currency{
		Denom: denom2,
		MetaData: &metadata.MetaData{
			Network:   metadata.Network_DEVNET,
			UpdatedAt: timestamppb.Now(),
			CreatedAt: timestamppb.Now(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	app := order.NewApplicationWithClients(currencyService, orderService)
	orders, err := app.OrderBookRelevantOrders(metadata.Network_DEVNET, "dextestdenom0-devcore1p0edzyzpazpt68vdrjy20c42lvwsjpvfzahygs", "udevcore", 20, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(orders)
}
