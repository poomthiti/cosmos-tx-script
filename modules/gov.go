package modules

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func GenerateGovMsgs(address sdk.AccAddress) (*govtypes.MsgSubmitProposal, *govtypes.MsgVote, *govtypes.MsgVoteWeighted, *govtypes.MsgDeposit) {
	content := govtypes.ContentFromProposalType("Test Proposal", "This proposal is for testing only.", "Text")
	msgSubmitProposal, _ := govtypes.NewMsgSubmitProposal(content, sdk.NewCoins(sdk.NewCoin("uosmo", sdk.NewInt(125000000))), address)

	return msgSubmitProposal,
		govtypes.NewMsgVote(address, 1, govtypes.OptionYes),
		govtypes.NewMsgVoteWeighted(
			address,
			1,
			govtypes.WeightedVoteOptions{
				govtypes.WeightedVoteOption{Option: govtypes.OptionYes, Weight: sdk.NewDec(70)},
				govtypes.WeightedVoteOption{Option: govtypes.OptionAbstain, Weight: sdk.NewDec(30)}}),
		govtypes.NewMsgDeposit(address, 1, sdk.NewCoins(sdk.NewCoin("uosmo", sdk.NewInt(100))))
}
