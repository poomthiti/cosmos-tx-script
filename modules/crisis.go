package modules

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

func GenerateCrisisMsgs(address sdk.AccAddress) *crisistypes.MsgVerifyInvariant {
	return crisistypes.NewMsgVerifyInvariant(address, "bank", "total-supply")
}
