package jsonapi

import (
	"context"
	"errors"
	"fmt"
)

type Cardholder struct {
	ID           uint
	CardholderID uint
	Cards        []KeepUnknown[Card]
	FirstName    string
	LastName     string

	/*
		{
		  "ID": 0,
		  "CardholderID": 0,
		  "CardholderTypeID": 0,
		  "Cards": [
		    {< see Card >}
		  ],
		  "City": "string",
		  "Companies": [
		    0
		  ],
		  "Country": "string",
		  "CustomFields": {
		    "additionalProp1": {},
		    "additionalProp2": {},
		    "additionalProp3": {}
		  },
		  "Department": "string",
		  "Department1": 0,
		  "Department2": 0,
		  "Email": "string",
		  "Extension": "string",
		  "FingerPrint": "string",
		  "FirstName": "string",
		  "Initials": "string",
		  "LastName": "string",
		  "Name": "string",
		  "Notes": "string",
		  "Phone": "string",
		  "Picture": "string",
		  "PictureSecond": "string",
		  "Postal": "string",
		  "Signature": "string",
		  "SignatureSecond": "string",
		  "State": "string",
		  "Street": "string",
		  "TemplateIndex": 0
		}
	*/
}

func (ch Cardholder) validateForWrite() error {
	if ch.FirstName == "" && ch.LastName == "" {
		return errors.New("empty name")
	}
	for i, cardK := range ch.Cards {
		card := cardK.Known
		if card.CardNumber == 0 {
			return fmt.Errorf("card[%d] missing card number", i)
		}
	}
	return nil
}

func (ch Cardholder) validateAfterWrite(written Cardholder) error {
	// compare cards, axiom will silently drop cards if they fail validation
	if len(ch.Cards) > len(written.Cards) {
		// use > in case the server cardholder already has some card holders
		return errors.New("cardholder localCard failed server validation and was not written")
	}
	writtenCardCount := 0
	for i, localCard := range ch.Cards {
		for _, serverCard := range written.Cards {
			if localCard.Known.CardNumber == serverCard.Known.CardNumber {
				if localCard.Known.AccessLevel != serverCard.Known.AccessLevel ||
					localCard.Known.CardType != serverCard.Known.CardType ||
					localCard.Known.ActiveDate != serverCard.Known.ActiveDate ||
					localCard.Known.ExpiryDate != serverCard.Known.ExpiryDate {
					return fmt.Errorf("card[%d] with card number %d is different on the server, suspected duplicate card number", i, localCard.Known.CardNumber)
				}
				writtenCardCount++
			}
		}
	}
	if writtenCardCount != len(ch.Cards) {
		return fmt.Errorf("only %d / %d cardholder cards were written to the server", writtenCardCount, len(ch.Cards))
	}
	return nil
}

func (c *Client) CreateCardholder(ctx context.Context, cardholder Cardholder) (KeepUnknown[Cardholder], error) {
	if err := cardholder.validateForWrite(); err != nil {
		return KeepUnknown[Cardholder]{}, err
	}

	newCardholder := KeepUnknown[Cardholder]{}
	if err := c.get(ctx, "/cardholder/new", &newCardholder); err != nil {
		return KeepUnknown[Cardholder]{}, err
	}

	newCardholder.Known.Cards = cardholder.Cards
	newCardholder.Known.FirstName = cardholder.FirstName
	newCardholder.Known.LastName = cardholder.LastName

	var id uint
	if err := c.post(ctx, "/cardholder/update", newCardholder, &id); err != nil {
		return KeepUnknown[Cardholder]{}, err
	}

	// Axiom will silently drop cards on the floor if they fail validation or have duplicated card numbers.
	// We read the card holder back to make sure we have everything and it matches
	writtenCardHolder := KeepUnknown[Cardholder]{}
	if err := c.post(ctx, "/cardholder/one", GetOneRequest{ID: id}, &writtenCardHolder); err != nil {
		return KeepUnknown[Cardholder]{}, err
	}
	if err := newCardholder.Known.validateAfterWrite(writtenCardHolder.Known); err != nil {
		return KeepUnknown[Cardholder]{}, err
	}

	return writtenCardHolder, nil
}
