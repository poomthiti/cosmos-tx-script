package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func FindAttrValue(events sdk.StringEvents, eventType string, attrKey string) string {
	for i := range events {
		if event := events[i]; event.GetType() == eventType {
			for idx := range event.GetAttributes() {
				if attribute := event.GetAttributes()[idx]; attribute.GetKey() == attrKey {
					return attribute.GetValue()
				}
			}
			return ""
		}
	}
	return ""
}
