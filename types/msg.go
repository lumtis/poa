package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// verify interface at compile time
var _ sdk.Msg = &MsgSubmitApplication{}

// var _ sdk.Msg = &MsgApproveApplication{}
// var _ sdk.Msg = &MsgRejectApplication{}

type MsgSubmitApplication struct {
	Candidate Validator `json:"validator"`
}

func NewMsgSubmitApplication(candidate Validator) MsgSubmitApplication {
	return MsgSubmitApplication{
		Candidate: candidate,
	}
}

const SubmitApplicationConst = "SubmitApplication"

func (msg MsgSubmitApplication) Route() string { return RouterKey }
func (msg MsgSubmitApplication) Type() string  { return SubmitApplicationConst }
func (msg MsgSubmitApplication) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Candidate.GetOperator())}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgSubmitApplication) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgSubmitApplication) ValidateBasic() error {
	return msg.Candidate.CheckValid()
}
