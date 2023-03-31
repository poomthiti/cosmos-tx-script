package main

import (
	"fmt"
	"os"
	"reflect"
	"scripts/cosmos-sdk/modules"
	"scripts/cosmos-sdk/utils"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/osmosis-labs/osmosis/v15/app"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

	// authz (3)
	msgGrant, msgRevoke, msgExec := modules.GenerateAuthzMsgs(address, granteeAddr)
	// bank (2)
	msgSend, msgMultiSend := modules.GenerateBankMsgs(address)
	// crisis (1)
	// msgVerifyInvariant := modules.GenerateCrisisMsgs(address)
	// distribution (4)
	_, msgWithdrawDelegatorReward, _, msgFundCommunityPool := modules.GenerateDistributionMsgs(address)
	// evidence (1)
	// msgSubmitEvidence := modules.GenerateEvidenceMsgs(address)
	// feegrant (2)
	// msgGrantAllowance, msgRevokeAllowance := modules.GenerateFeeGrantMsgs(address, granteeAddr)
	// gov (4)
	_, msgVote, msgVoteWeighted, msgDeposit := modules.GenerateGovMsgs(address)
	// slashing (1)
	msgUnjail := modules.GenerateSlashingMsgs()
	// staking (5)
	_, _, _, msgBeginRedelegate, msgUndelegate := modules.GenerateStakingMsgs(address)

	msgs := modules.Msgs{
		// MsgSubmitProposal:     msgSubmitProposal,
		// MsgCreateValidator:    msgCreateValidator,
		// MsgEditValidator:      msgEditValidator,
		// MsgDelegate:           msgDelegate,
		// MsgSetWithdrawAddress: msgSetWithdrawAddress,
		MsgGrant:             msgGrant,
		MsgExec:              msgExec,
		MsgRevoke:            msgRevoke,
		MsgSend:              msgSend,
		MsgMultiSend:         msgMultiSend,
		MsgFundCommunityPool: msgFundCommunityPool,
		// ------------Not possible-----------------
		// MsgVerifyInvariant:         msgVerifyInvariant,
		// ------------Not possible-----------------
		// MsgWithdrawValidatorCommission: msgWithdrawValidatorCommission,
		// ------------Not possible-----------------
		// MsgSubmitEvidence: msgSubmitEvidence,
		// ------x/feegrant is not registered in Osmosis chain-------
		// MsgGrantAllowance:  msgGrantAllowance,
		// MsgRevokeAllowance: msgRevokeAllowance,
		// ------Entering proposal voting period costs a lot-------
		MsgVote:                    msgVote,
		MsgVoteWeighted:            msgVoteWeighted,
		MsgDeposit:                 msgDeposit,
		MsgUnjail:                  msgUnjail,
		MsgWithdrawDelegatorReward: msgWithdrawDelegatorReward,
		MsgBeginRedelegate:         msgBeginRedelegate,
		MsgUndelegate:              msgUndelegate,
	}

	sequence := acc.GetSequence()
	values := reflect.ValueOf(msgs)
	types := values.Type()
	var proposalId string

	// ------------------------------------------------
	// ----------------SUCCESS CASES-------------------
	// ------------------------------------------------
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

		_, simGasUsed, err := tx.CalculateGas(clientCtx.GRPCClient, txf.WithSequence(sequence), msg)
		if err != nil {
			panic(err)
		}

		txf = txf.WithGas(simGasUsed).WithSequence(sequence)

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

		sequence += 1
		fmt.Printf("[SUCCESS]Tx broadcast successful. Height: %d, RawLog: %s,TxHash: %s\n", res.Height, res.RawLog, res.TxHash)
		println("---------------------------------------------------------------------------------------------------------------------------------")
	}

	// ------------------------------------------------
	// ----------------FAILED CASES--------------------
	// ------------------------------------------------
	dummyValAddr := "osmovaloper1l0ta4rw7zauqplzhsvcsgxveuqptauf6e4eg7a"
	for i := 0; i < values.NumField(); i++ {
		msgName := types.Field(i).Name
		var msg sdk.Msg
		println(msgName)

		switch msgName {
		case "MsgSubmitProposal":
			typedMsg := values.Field(i).Interface().(*govtypes.MsgSubmitProposal)
			typedMsg.InitialDeposit = sdk.NewCoins(sdk.NewCoin("uosmo", sdk.NewInt(1000000)))
			msg = typedMsg
		case "MsgCreateValidator":
			typedMsg := values.Field(i).Interface().(*stakingtypes.MsgCreateValidator)
			msg = typedMsg
		case "MsgEditValidator":
			typedMsg := values.Field(i).Interface().(*stakingtypes.MsgEditValidator)
			minDelegate := sdk.NewInt(1000000)
			typedMsg.MinSelfDelegation = &minDelegate
			msg = typedMsg
		case "MsgDelegate":
			typedMsg := values.Field(i).Interface().(*stakingtypes.MsgDelegate)
			typedMsg.Amount = sdk.NewCoin("test", sdk.NewInt(1000000000000000000))
			msg = typedMsg
		case "MsgSetWithdrawAddress":
			typedMsg := values.Field(i).Interface().(*distributiontypes.MsgSetWithdrawAddress)
			typedMsg.WithdrawAddress = "FAILED"
			msg = typedMsg
		case "MsgGrant":
			typedMsg := values.Field(i).Interface().(*authztypes.MsgGrant)
			typedMsg.Grant.Expiration = time.Date(2022, time.August, 22, 11, 11, 11, 11, time.UTC)
			msg = typedMsg
		case "MsgExec":
			failedSend := msgSend
			failedSend.Amount = (sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(10000000000))))
			failedMsgExec := authztypes.NewMsgExec(address, []sdk.Msg{failedSend})
			msg = &failedMsgExec
		case "MsgRevoke":
			typedMsg := values.Field(i).Interface().(*authztypes.MsgRevoke)
			typedMsg.MsgTypeUrl = "FAILED"
			msg = typedMsg
		case "MsgSend":
			typedMsg := values.Field(i).Interface().(*banktypes.MsgSend)
			typedMsg.Amount = sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(10000000000)))
			msg = typedMsg
		case "MsgMultiSend":
			typedMsg := values.Field(i).Interface().(*banktypes.MsgMultiSend)
			coin := sdk.Coin{Denom: "ufailed", Amount: sdk.NewInt(1)}
			amount := sdk.NewCoins(coin)
			typedMsg.Inputs = []banktypes.Input{{Address: address.String(), Coins: amount}}
			typedMsg.Outputs = []banktypes.Output{{Address: address.String(), Coins: amount}}
			msg = typedMsg
		case "MsgFundCommunityPool":
			typedMsg := values.Field(i).Interface().(*distributiontypes.MsgFundCommunityPool)
			typedMsg.Amount = sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(10000000000)))
			msg = typedMsg
		case "MsgVote":
			typedMsg := values.Field(i).Interface().(*govtypes.MsgVote)
			msg = typedMsg
		case "MsgVoteWeighted":
			typedMsg := values.Field(i).Interface().(*govtypes.MsgVoteWeighted)
			msg = typedMsg
		case "MsgDeposit":
			typedMsg := values.Field(i).Interface().(*govtypes.MsgDeposit)
			msg = typedMsg
		case "MsgUnjail":
			typedMsg := values.Field(i).Interface().(*slashingtypes.MsgUnjail)
			msg = typedMsg
		case "MsgWithdrawValidatorCommission":
			typedMsg := values.Field(i).Interface().(*distributiontypes.MsgWithdrawValidatorCommission)
			typedMsg.ValidatorAddress = dummyValAddr
			msg = typedMsg
		case "MsgWithdrawDelegatorReward":
			typedMsg := values.Field(i).Interface().(*distributiontypes.MsgWithdrawDelegatorReward)
			typedMsg.ValidatorAddress = dummyValAddr
			msg = typedMsg
		case "MsgBeginRedelegate":
			typedMsg := values.Field(i).Interface().(*stakingtypes.MsgBeginRedelegate)
			typedMsg.ValidatorSrcAddress = dummyValAddr
			msg = typedMsg
		case "MsgUndelegate":
			typedMsg := values.Field(i).Interface().(*stakingtypes.MsgUndelegate)
			typedMsg.ValidatorAddress = dummyValAddr
			msg = typedMsg
		}

		txf = txf.WithGas(1_00_000).WithSequence(sequence)

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

		fmt.Printf("[FAILED] Tx broadcast successful. Height: %d, RawLog: %s,TxHash: %s\n", res.Height, res.RawLog, res.TxHash)
		println("---------------------------------------------------------------------------------------------------------------------------------")
		sequence += 1
	}
}
