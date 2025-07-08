package baremetal

import (
	"context"
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var offerAvailabilityMarshalSpecs = human.EnumMarshalSpecs{
	baremetal.OfferStockEmpty:     &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "empty"},
	baremetal.OfferStockLow:       &human.EnumMarshalSpec{Attribute: color.FgYellow, Value: "low"},
	baremetal.OfferStockAvailable: &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "available"},
}

func listOfferMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp baremetal.Offer
	baremetalOffer := tmp(i.(baremetal.Offer))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Disks",
			Title:     "Disks",
		},
		{
			FieldName: "CPUs",
			Title:     "CPUs",
		},
		{
			FieldName: "Memories",
			Title:     "Memories",
		},
		{
			FieldName: "Options",
			Title:     "Options",
		},
		{
			FieldName: "Bandwidth",
			Title:     "Bandwidth(Mbit/s)",
		},
		{
			FieldName: "PrivateBandwidth",
			Title:     "PrivateBandwidth(Mbit/s)",
		},
	}
	baremetalOffer.PrivateBandwidth /= 1000000
	baremetalOffer.Bandwidth /= 1000000
	str, err := human.Marshal(baremetalOffer, opt)
	if err != nil {
		return "", err
	}
	return str, nil
}

func serverOfferListBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		req := argsI.(*baremetal.ListOffersRequest)
		baremetalAPI := baremetal.NewAPI(core.ExtractClient(ctx))

		offers, err := baremetalAPI.ListOffers(req, scw.WithAllPages())
		if err != nil {
			return nil, err
		}

		return offers.Offers, nil
	}

	return c
}
