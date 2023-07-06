package flixkit

import (
	"reflect"
	"testing"
)

var template = `{
	  "f_type": "InteractionTemplate",
	  "f_version": "1.0.0",
	  "id": "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	  "data": {
		"type": "transaction",
		"interface": "",
		"messages": {
		  "title": {
			"i18n": {
			  "en-US": "Transfer Tokens"
			}
		  },
		  "description": {
			"i18n": {
			  "en-US": "Transfer tokens from one account to another"
			}
		  }
		},
		"cadence": "import FungibleToken from 0xFUNGIBLETOKENADDRESS\ntransaction(amount: UFix64, to: Address) {\nlet vault: @FungibleToken.Vault\nprepare(signer: AuthAccount) {\nself.vault <- signer\n.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n.withdraw(amount: amount)\n}\nexecute {\ngetAccount(to)\n.getCapability(/public/flowTokenReceiver)!\n.borrow<&{FungibleToken.Receiver}>()!\n.deposit(from: <-self.vault)\n}\n}",
		"dependencies": {
		  "0xFUNGIBLETOKENADDRESS": {
			"FungibleToken": {
			  "mainnet": {
				"address": "0xf233dcee88fe0abe",
				"fq_address": "A.0xf233dcee88fe0abe.FungibleToken",
				"contract": "FungibleToken",
				"pin": "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				"pin_block_height": 34166296
			  },
			  "testnet": {
				"address": "0x9a0766d93b6608b7",
				"fq_address": "A.0x9a0766d93b6608b7.FungibleToken",
				"contract": "FungibleToken",
				"pin": "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				"pin_block_height": 74776482
			  }
			}
		  }
		},
		"arguments": {
		  "amount": {
			"index": 0,
			"type": "UFix64",
			"messages": {
			  "title": {
				"i18n": {
				  "en-US": "The amount of FLOW tokens to send"
				}
			  }
			},
			"balance": ""
		  },
		  "to": {
			"index": 1,
			"type": "Address",
			"messages": {
			  "title": {
				"i18n": {
				  "en-US": "The Flow account the tokens will go to"
				}
			  }
			},
			"balance": ""
		  }
		}
	  }
	}`

var parsedTemplate = &FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: Data{
		Type:      "transaction",
		Interface: "",
		Messages: Messages{
			Title: &Title{
				I18N: map[string]string{
					"en-US": "Transfer Tokens",
				},
			},
			Description: &Description{
				I18N: map[string]string{
					"en-US": "Transfer tokens from one account to another",
				},
			},
		},
		Cadence: "import FungibleToken from 0xFUNGIBLETOKENADDRESS\ntransaction(amount: UFix64, to: Address) {\nlet vault: @FungibleToken.Vault\nprepare(signer: AuthAccount) {\nself.vault <- signer\n.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n.withdraw(amount: amount)\n}\nexecute {\ngetAccount(to)\n.getCapability(/public/flowTokenReceiver)!\n.borrow<&{FungibleToken.Receiver}>()!\n.deposit(from: <-self.vault)\n}\n}",
		Dependencies: Dependencies{
			"0xFUNGIBLETOKENADDRESS": Contracts{
				"FungibleToken": Networks{
					"mainnet": Network{
						Address:        "0xf233dcee88fe0abe",
						FqAddress:      "A.0xf233dcee88fe0abe.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 34166296,
					},
					"testnet": Network{
						Address:        "0x9a0766d93b6608b7",
						FqAddress:      "A.0x9a0766d93b6608b7.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 74776482,
					},
				},
			},
		},
		Arguments: Arguments{
			"amount": Argument{
				Index: 0,
				Type:  "UFix64",
				Messages: Messages{
					Title: &Title{
						I18N: map[string]string{
							"en-US": "The amount of FLOW tokens to send",
						},
					},
				},
				Balance: "",
			},
			"to": Argument{
				Index: 1,
				Type:  "Address",
				Messages: Messages{
					Title: &Title{
						I18N: map[string]string{
							"en-US": "The Flow account the tokens will go to",
						},
					},
				},
				Balance: "",
			},
		},
	},
}

func TestParseTemplate(t *testing.T) {
	parsedTemplate, err := ParseTemplate(template)
	if err != nil {
		t.Fatalf("ParseTemplate() err = %v; want nil", err)
	}
	if !reflect.DeepEqual(parsedTemplate, parsedTemplate) {
		t.Errorf("ParseTemplate() = %v; want %v", parsedTemplate, parsedTemplate)
	}
}
