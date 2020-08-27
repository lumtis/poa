package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// verify interface at compile time
var _ sdk.Msg = &MsgSubmitApplication{}
var _ sdk.Msg = &MsgVote{}

/**
 * MsgSubmitApplication
 */
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

/**
 * MsgVote
 */

type MsgVote struct {
	VoteType      uint16         `json:"type"`
	VoterAddr     sdk.ValAddress `json:"voter"`
	CandidateAddr sdk.ValAddress `json:"candidate"`
	Approve       bool           `json:"approve"`
}

func NewMsgVote(voteType uint16, voter sdk.ValAddress, candidate sdk.ValAddress, approve bool) MsgVote {
	return MsgVote{
		VoteType:      voteType,
		VoterAddr:     voter,
		CandidateAddr: candidate,
		Approve:       approve,
	}
}

const VoteConst = "Vote"

const (
	VoteTypeApplication  uint16 = iota
	VoteTypeKickProposal uint16 = iota
)

func (msg MsgVote) Route() string { return RouterKey }
func (msg MsgVote) Type() string  { return VoteConst }
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.VoterAddr)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgVote) ValidateBasic() error {
	if msg.VoterAddr.Empty() || msg.CandidateAddr.Empty() {
		return sdkerrors.Wrap(ErrInvalidVoteMsg, "missing address")
	}
	if msg.VoteType != VoteTypeApplication && msg.VoteType != VoteTypeKickProposal {
		return sdkerrors.Wrap(ErrInvalidVoteMsg, "vote type incorrect")
	}

	return nil
}
