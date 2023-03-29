package modules

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

func GenerateSlashingMsgs() *slashingtypes.MsgUnjail {
	valAddress, _ := sdk.ValAddressFromBech32("osmovaloper1l0ta4rw7zauqplzhsvcsgxveuqptauf6e4eg7a")
	msgUnjail := slashingtypes.NewMsgUnjail(valAddress)
	return msgUnjail
}
