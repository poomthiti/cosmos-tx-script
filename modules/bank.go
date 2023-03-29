package modules

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func GenerateBankMsgs(addr sdk.AccAddress) (*banktypes.MsgSend, *banktypes.MsgMultiSend) {
	coin := sdk.Coin{Denom: "uosmo", Amount: sdk.NewInt(1)}
	amount := sdk.NewCoins(coin)
	input := []banktypes.Input{{Address: addr.String(), Coins: amount}}
	output := []banktypes.Output{{Address: addr.String(), Coins: amount}}

	return banktypes.NewMsgSend(addr, addr, amount), banktypes.NewMsgMultiSend(input, output)
}
