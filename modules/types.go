package modules

import (
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Msgs struct {
	MsgSubmitProposal     *govtypes.MsgSubmitProposal
	MsgCreateValidator    *stakingtypes.MsgCreateValidator
	MsgEditValidator      *stakingtypes.MsgEditValidator
	MsgDelegate           *stakingtypes.MsgDelegate
	MsgSetWithdrawAddress *distributiontypes.MsgSetWithdrawAddress
	MsgGrant              *authztypes.MsgGrant
	MsgExec               *authztypes.MsgExec
	MsgRevoke             *authztypes.MsgRevoke
	MsgSend               *banktypes.MsgSend
	MsgMultiSend          *banktypes.MsgMultiSend
	MsgFundCommunityPool  *distributiontypes.MsgFundCommunityPool
	// MsgVerifyInvariant             *crisistypes.MsgVerifyInvariant
	// MsgSubmitEvidence              *evidencetypes.MsgSubmitEvidence
	// MsgGrantAllowance              *feegranttypes.MsgGrantAllowance
	// MsgRevokeAllowance             *feegranttypes.MsgRevokeAllowance
	MsgVote                        *govtypes.MsgVote
	MsgVoteWeighted                *govtypes.MsgVoteWeighted
	MsgDeposit                     *govtypes.MsgDeposit
	MsgUnjail                      *slashingtypes.MsgUnjail
	MsgWithdrawValidatorCommission *distributiontypes.MsgWithdrawValidatorCommission
	MsgWithdrawDelegatorReward     *distributiontypes.MsgWithdrawDelegatorReward
	MsgBeginRedelegate             *stakingtypes.MsgBeginRedelegate
	MsgUndelegate                  *stakingtypes.MsgUndelegate
}
