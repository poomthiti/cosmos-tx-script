package main

import (
	"fmt"
	"os"
	"reflect"
	"scripts/cosmos-sdk/modules"
	"scripts/cosmos-sdk/utils"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/osmosis-labs/osmosis/v15/app"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	sdk.GetConfig().SetBech32PrefixForAccount(os.Getenv("PREFIX"), "")
	sendTx()
}

func sendTx() {
	MNEMONIC := os.Getenv("MNEMONIC")
	RPC_PATH := os.Getenv("RPC_PATH")
	CHAIN_ID := os.Getenv("CHAIN_ID")
	GAS_PRICES := os.Getenv("GAS_PRICES")
	GRPC_PATH := os.Getenv("GRPC_PATH")

	// Setup keyring
	kb := keyring.NewInMemory()
	path := sdk.GetConfig().GetFullBIP44Path()
	key, err := kb.NewAccount("sdk-ja", MNEMONIC, "", path, hd.Secp256k1)
	if err != nil {
		panic(err)
	}
	println("Account:", key.GetAddress().String())

	address := key.GetAddress()
	key.GetPubKey()

	// Create client
	clientNode, err := client.NewClientFromNode(RPC_PATH)
	if err != nil {
		panic(err)
	}

	grpcClient, err := grpc.Dial(GRPC_PATH, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	clientCtx := client.Context{
		Client:            clientNode,
		GRPCClient:        grpcClient,
		ChainID:           CHAIN_ID,
		NodeURI:           RPC_PATH,
		InterfaceRegistry: app.MakeEncodingConfig().InterfaceRegistry,
		TxConfig:          app.MakeEncodingConfig().TxConfig,
		Keyring:           kb,
	}

	// Retrieve account info
	accountRetriever := authtypes.AccountRetriever{}
	acc, err := accountRetriever.GetAccount(clientCtx, key.GetAddress())
	if err != nil {
		panic(err)
	}
	granteeAddr, _ := sdk.AccAddressFromBech32("osmo1wke7j8f5kgnnacs3avchcj6fvvdtvrsalzmddx")
	grantee, err := accountRetriever.GetAccount(clientCtx, granteeAddr)
	if err != nil {
		panic(err)
	}

	// Create transaction factory
	txf := tx.Factory{}.
		WithKeybase(kb).
		WithTxConfig(app.MakeEncodingConfig().TxConfig).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithAccountNumber(acc.GetAccountNumber()).
		WithGasPrices(GAS_PRICES).
		WithGasAdjustment(2).
		WithChainID(CHAIN_ID).
		WithMemo("").
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	// granteeTxf := tx.Factory{}.
	// 	WithKeybase(kb).
	// 	WithTxConfig(app.MakeEncodingConfig().TxConfig).
	// 	WithAccountRetriever(clientCtx.AccountRetriever).
	// 	WithAccountNumber(grantee.GetAccountNumber()).
	// 	WithSequence(grantee.GetSequence()).
	// 	WithGasPrices(GAS_PRICES).
	// 	WithGasAdjustment(2).
	// 	WithChainID(CHAIN_ID).
	// 	WithMemo("").
	// 	WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	// authz (3)
	msgGrant, msgRevoke, msgExec := modules.GenerateAuthzMsgs(address, granteeAddr)
	// bank (2)
	msgSend, msgMultiSend := modules.GenerateBankMsgs(address)
	// crisis (1)
	// msgVerifyInvariant := modules.GenerateCrisisMsgs(address)
	// distribution (4)
	msgSetWithdrawAddress, msgWithdrawDelegatorReward, _, msgFundCommunityPool := modules.GenerateDistributionMsgs(address)
	// evidence (1)
	// msgSubmitEvidence := modules.GenerateEvidenceMsgs(address)
	// feegrant (2)
	// msgGrantAllowance, msgRevokeAllowance := modules.GenerateFeeGrantMsgs(address, granteeAddr)
	// gov (4)
	msgSubmitProposal, _, _, msgDeposit := modules.GenerateGovMsgs(address)
	// slashing (1)
	// msgUnjail := modules.GenerateSlashingMsgs()
	// staking (5)
	_, msgEditValidator, msgDelegate, msgBeginRedelegate, msgUndelegate := modules.GenerateStakingMsgs(address)

	msgs := modules.Msgs{
		MsgSubmitProposal: msgSubmitProposal,
		// MsgCreateValidator:    msgCreateValidator,
		MsgEditValidator:      msgEditValidator,
		MsgDelegate:           msgDelegate,
		MsgSetWithdrawAddress: msgSetWithdrawAddress,
		MsgGrant:              msgGrant,
		MsgExec:               msgExec,
		MsgRevoke:             msgRevoke,
		MsgSend:               msgSend,
		MsgMultiSend:          msgMultiSend,
		MsgFundCommunityPool:  msgFundCommunityPool,
		// ------------Not possible-----------------
		// MsgVerifyInvariant:         msgVerifyInvariant,
		// ------x/feegrant is not registered in Osmosis chain-------
		// MsgWithdrawValidatorCommission: msgWithdrawValidatorCommission,
		// ------------Not possible-----------------
		// MsgSubmitEvidence: msgSubmitEvidence,
		// ------x/feegrant is not registered in Osmosis chain-------
		// MsgGrantAllowance:  msgGrantAllowance,
		// MsgRevokeAllowance: msgRevokeAllowance,
		// ------Entering proposal voting period costs a lot-------
		// MsgVote:         msgVote,
		// MsgVoteWeighted: msgVoteWeighted,
		MsgDeposit: msgDeposit,
		// ------Need to get your own validator jailed-------
		// MsgUnjail:          msgUnjail,
		MsgWithdrawDelegatorReward: msgWithdrawDelegatorReward,
		MsgBeginRedelegate:         msgBeginRedelegate,
		MsgUndelegate:              msgUndelegate,
	}

	sequence := acc.GetSequence()
	values := reflect.ValueOf(msgs)
	types := values.Type()
	var proposalId string

	for i := 0; i < values.NumField(); i++ {
		msgName := types.Field(i).Name
		fmt.Println(msgName)
		var msg sdk.Msg
		switch msgName {
		case "MsgVote":
			typedMsg := values.Field(i).Interface().(*govtypes.MsgVote)
			id, err := strconv.ParseUint(proposalId, 10, 64)
			if err != nil {
				panic("Failed parsing proposalId")
			}
			typedMsg.ProposalId = id
			msg = typedMsg
		case "MsgVoteWeighted":
			typedMsg := values.Field(i).Interface().(*govtypes.MsgVoteWeighted)
			id, err := strconv.ParseUint(proposalId, 10, 64)
			if err != nil {
				panic("Failed parsing proposalId")
			}
			typedMsg.ProposalId = id
			msg = typedMsg
		case "MsgDeposit":
			typedMsg := values.Field(i).Interface().(*govtypes.MsgDeposit)
			id, err := strconv.ParseUint(proposalId, 10, 64)
			if err != nil {
				panic("Failed parsing proposalId")
			}
			typedMsg.ProposalId = id
			msg = typedMsg
		default:
			msg = values.Field(i).Interface().(sdk.Msg)
		}

		if msgName == "MsgExec" {
			txf = txf.WithSequence(grantee.GetSequence())
		} else {
			txf = txf.WithSequence(sequence)
			sequence += 1
		}
		_, simGasUsed, err := tx.CalculateGas(clientCtx.GRPCClient, txf, msg)
		if err != nil {
			panic(err)
		}
		fmt.Printf("SIM GAS USED: %s\n", strconv.FormatUint(simGasUsed, 10))

		txf = txf.WithGas(simGasUsed)

		txb, err := tx.BuildUnsignedTx(txf, msg)
		if err != nil {
			panic(err)
		}

		err = tx.Sign(txf, key.GetName(), txb, true)
		if err != nil {
			panic(err)
		}

		txBytes, err := clientCtx.TxConfig.TxEncoder()(txb.GetTx())
		if err != nil {
			panic(err)
		}

		res, err := clientCtx.BroadcastTxCommit(txBytes)
		if err != nil {
			panic(err)
		}

		if msgName == "MsgSubmitProposal" {
			proposalId = utils.FindAttrValue(res.Logs[0].GetEvents(), "submit_proposal", "proposal_id")
		}

		fmt.Printf("Tx broadcast successful. Message: %s, TxHash: %s\n", sdk.MsgTypeURL(msg), res.TxHash)
	}
}
