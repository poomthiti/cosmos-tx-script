package modules

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
)

func GenerateFeeGrantMsgs(address sdk.AccAddress, granteeAddr sdk.AccAddress) (*feegranttypes.MsgGrantAllowance, *feegranttypes.MsgRevokeAllowance) {
	expiration := time.Now().AddDate(1, 0, 0)
	msgGrantAllowance, _ := feegranttypes.NewMsgGrantAllowance(&feegranttypes.BasicAllowance{
		SpendLimit: sdk.NewCoins(sdk.NewCoin("uosmo", sdk.NewInt(100))),
		Expiration: &expiration,
	}, address, granteeAddr)
	msgRevokeAllowance := feegranttypes.NewMsgRevokeAllowance(address, granteeAddr)
	return msgGrantAllowance, &msgRevokeAllowance
}
