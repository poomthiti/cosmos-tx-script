package modules

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func GenerateDistributionMsgs(address sdk.AccAddress) (
	*distributiontypes.MsgSetWithdrawAddress,
	*distributiontypes.MsgWithdrawDelegatorReward,
	*distributiontypes.MsgWithdrawValidatorCommission,
	*distributiontypes.MsgFundCommunityPool) {
	// Hardcode?
	dstAddr, _ := sdk.AccAddressFromBech32("osmo1c584m4lq25h83yp6ag8hh4htjr92d954vklzja")
	dstValAddr := sdk.ValAddress(dstAddr)
	amount := sdk.NewCoins(sdk.NewCoin("uosmo", sdk.NewInt(100)))

	return distributiontypes.NewMsgSetWithdrawAddress(address, address),
		distributiontypes.NewMsgWithdrawDelegatorReward(address, dstValAddr),
		distributiontypes.NewMsgWithdrawValidatorCommission(dstValAddr),
		distributiontypes.NewMsgFundCommunityPool(amount, address)
}
