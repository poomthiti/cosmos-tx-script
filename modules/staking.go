package modules

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func GenerateStakingMsgs(address sdk.AccAddress) (
	*stakingtypes.MsgCreateValidator,
	*stakingtypes.MsgEditValidator,
	*stakingtypes.MsgDelegate,
	*stakingtypes.MsgBeginRedelegate,
	*stakingtypes.MsgUndelegate) {
	srcVal := sdk.ValAddress(address)
	pubkey := ed25519.GenPrivKey().PubKey()
	dstAddr, _ := sdk.AccAddressFromBech32("osmo1c584m4lq25h83yp6ag8hh4htjr92d954vklzja")
	dstValAddr := sdk.ValAddress(dstAddr)
	// newRates := sdk.NewDecWithPrec(11, 2)
	// newMinSelfDelegation := sdk.NewInt(1100000)
	msgCreateValidator, _ := stakingtypes.NewMsgCreateValidator(srcVal, pubkey, sdk.NewCoin("uosmo",
		sdk.NewInt(1000000)),
		stakingtypes.NewDescription("TestVal", "", "", "", "TestValidator - For testing only"),
		stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(1, 1), sdk.NewDecWithPrec(3, 1), sdk.NewDecWithPrec(1, 2)),
		sdk.NewInt(1000000))

	return msgCreateValidator,
		stakingtypes.NewMsgEditValidator(
			srcVal,
			stakingtypes.NewDescription("TestValEdit", "", "", "", "TestValidatorEdit - For testing edit only"),
			nil,
			nil),
		stakingtypes.NewMsgDelegate(address, dstValAddr, sdk.NewCoin("uosmo", sdk.NewInt(1000))),
		stakingtypes.NewMsgBeginRedelegate(address, dstValAddr, srcVal, sdk.NewCoin("uosmo", sdk.NewInt(500))),
		stakingtypes.NewMsgUndelegate(address, srcVal, sdk.NewCoin("uosmo", sdk.NewInt(200)))
}
