package modules

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

func GenerateSlashingMsgs() *slashingtypes.MsgUnjail {
	valAddress, _ := sdk.ValAddressFromBech32("osmovaloper1acqpnvg2t4wmqfdv8hq47d3petfksjs5ejrkrx")
	msgUnjail := slashingtypes.NewMsgUnjail(valAddress)
	return msgUnjail
}
