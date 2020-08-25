package cli

import (
	"bufio"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ltacker/poa/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	poaTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	poaTxCmd.AddCommand(flags.PostCommands(
		GetCmdSubmitApplication(cdc),
	)...)

	return poaTxCmd
}

// GetCmdSubmitApplication sends a new application to become a validator
func GetCmdSubmitApplication(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply to become a new validator in the network",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// Operator address is the sender
			accAddress := cliCtx.GetFromAddress()
			if accAddress.Empty() {
				return fmt.Errorf("Account address empty")
			}

			opAddress := sdk.ValAddress(accAddress)

			// Consensus public key for the validator
			pkStr, err := cmd.Flags().GetString(FlagPubKey)
			if err != nil {
				return fmt.Errorf("Cannot get pubkey flag: %v", err)
			}
			pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pkStr)
			if err != nil {
				return fmt.Errorf("Cannot convert pubkey: %v", err)
			}

			// Description of the candidate
			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			identity, _ := cmd.Flags().GetString(FlagIdentity)
			website, _ := cmd.Flags().GetString(FlagWebsite)
			security, _ := cmd.Flags().GetString(FlagSecurityContact)
			details, _ := cmd.Flags().GetString(FlagDetails)
			description := types.NewDescription(moniker, identity, website, security, details)

			candidateValidator := types.NewValidator(opAddress, pk, description)

			msg := types.NewMsgSubmitApplication(candidateValidator)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FlagSetDescriptionCreate())
	_ = cmd.MarkFlagRequired(FlagPubKey)

	return cmd
}
