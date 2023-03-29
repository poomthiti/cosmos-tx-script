package modules

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
)

func GenerateAuthzMsgs(address sdk.AccAddress, granteeAddr sdk.AccAddress) (*authztypes.MsgGrant, *authztypes.MsgRevoke, *authztypes.MsgExec) {
	typeUrl := "/cosmos.bank.v1beta1.MsgSend"
	authorization := authztypes.NewGenericAuthorization(typeUrl)
	expiry := time.Date(2023, time.August, 22, 11, 11, 11, 11, time.UTC)
	msgSend, _ := GenerateBankMsgs(address)

	msgGrant, _ := authztypes.NewMsgGrant(address, granteeAddr, authorization, expiry)
	msgRevoke := authztypes.NewMsgRevoke(address, granteeAddr, typeUrl)
	msgExec := authztypes.NewMsgExec(granteeAddr, []sdk.Msg{msgSend})

	return msgGrant, &msgRevoke, &msgExec
}
