package modules

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func GenerateEvidenceMsgs(address sdk.AccAddress) *evidencetypes.MsgSubmitEvidence {
	tmEvidence := abci.Evidence{
		Type: abci.EvidenceType_DUPLICATE_VOTE,
		Validator: abci.Validator{
			Address: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Power:   100,
		},
		Height:           1,
		Time:             time.Now(),
		TotalVotingPower: 100,
	}
	evidence := evidencetypes.FromABCIEvidence(tmEvidence).(*evidencetypes.Equivocation)
	msgSubmitEvidence, _ := evidencetypes.NewMsgSubmitEvidence(address, evidence)
	return msgSubmitEvidence
}
