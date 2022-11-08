package jsonapi

import (
	"context"
	"fmt"
)

// Card represents a "card object" as defined by Section 5 of the API documentation.
type Card struct {
	ID           int `json:"ID"`
	CardholderID int `json:"CardholderID"`
	CardNumber   int `json:"CardNumber"`
	CardType     int `json:"CardType"`

	/*
		{
		  "AccessLevel": 0,
		  "CardholderId": 0,
		  "AccessLevelSpecial": [
		    {
		      "ID": 0,
		      "ALID": 0,
		      "ALName": "string",
		      "APID": 0,
		      "Name": "string",
		      "APName": "string",
		      "TGID": 0,
		      "SALExpiryDate": "2021-07-13T16:12:21.263Z"
		    }
		  ],
		  "ActiveDate": "2021-07-13T16:12:21.263Z",
		  "AssetHolderNumber": 0,
		  "AutoVoidDate": "2021-07-13T16:12:21.263Z",
		  "ID": 0,
		  "CardNumber": 0,
		  "CardType": 0,
		  "DownloadPendingCommands": true,
		  "EscortRequired": true,
		  "ExpiryDate": "2021-07-13T16:12:21.263Z",
		  "ExtendedUnlockTime": true,
		  "FingerPrints": {
		    "CardNumber": 0,
		    "FingerIndex": 0,
		    "MasterFinger": true,
		    "Name": "string",
		    "Pin": 0,
		    "SiteCode": 0,
		    "Status": 0,
		    "Template1": "string",
		    "Template2": "string"
		  },
		  "IgnoreAPB": true,
		  "IgnoreAutovoid": true,
		  "IgnoreHighSecurity": true,
		  "IrisData": "string",
		  "IssueLevel": 0,
		  "Links": [
		    {
		      "NetworkID": 0,
		      "APID": 0,
		      "AccessPointName": "string",
		      "LinkID": 0,
		      "LinkName": "string",
		      "ID": 0,
		      "Name": "string"
		    }
		  ],
		  "MultiAccessLevel": [
		    {
		      "MultiAccessLevel": 0,
		      "MultiAccessLevelId": "string",
		      "MultiAccessLevelName": "string",
		      "MALExpiryDate": "2021-07-13T16:12:21.263Z",
		      "ID": 0,
		      "Name": "string"
		    }
		  ],
		  "Name": "string",
		  "Options": 0,
		  "OriginalVacations": {
		    "EndDate": "2021-07-13T16:12:21.263Z",
		    "StartDate": "2021-07-13T16:12:21.263Z",
		    "ID": 0,
		    "Name": "string"
		  },
		  "PIN": "string",
		  "PendingCommandChanged": true,
		  "ReaderAccess": [
		    {
		      "HSAccess": true,
		      "IsAccessLevelMember": true,
		      "LUAccess": true,
		      "ID": 0,
		      "Name": "string"
		    }
		  ],
		  "Status": 0,
		  "StealthModeSchedule": 0,
		  "TraceCard": true,
		  "UseCount": 0,
		  "Vacations": [
		    {
		      "EndDate": "2021-07-13T16:12:21.263Z",
		      "StartDate": "2021-07-13T16:12:21.263Z",
		      "ID": 0,
		      "Name": "string"
		    }
		  ]
		}


	*/
}

// CreateCard creates a new card with the given card information.
// Card.ID will be ignored if present in the given card.
func (c *Client) CreateCard(ctx context.Context, card Card) (KeepUnknown[Card], error) {
	var newCard KeepUnknown[Card]
	if err := c.get(ctx, "/card/new", &newCard); err != nil {
		return KeepUnknown[Card]{}, err
	}

	newCard.Known.CardholderID = card.CardholderID
	newCard.Known.CardNumber = card.CardNumber
	newCard.Known.CardType = card.CardType

	var updated bool
	if err := c.post(ctx, "/card/update", newCard, &updated); err != nil {
		return KeepUnknown[Card]{}, err
	}
	if !updated {
		return newCard, fmt.Errorf("/card/update returned false, not true")
	}
	return newCard, nil
}
