package v18de63

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/tendermint/tendermint/crypto"
	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ------------------------------------------
//				Keys
// ------------------------------------------

const ModuleName = "auth"

// ------------------------------------------
//				GenesisState
// ------------------------------------------

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	Params   Params                   `json:"params" yaml:"params"`
	Accounts GenesisAccounts `json:"accounts" yaml:"accounts"`
}
// TODO copy in methods for validation?

// ------------------------------------------
//				Params
// ------------------------------------------
// TODO import from somewhere if they haven't changed?
// Params defines the parameters for the auth module.
type Params struct {
	MaxMemoCharacters      uint64 `json:"max_memo_characters" yaml:"max_memo_characters"`
	TxSigLimit             uint64 `json:"tx_sig_limit" yaml:"tx_sig_limit"`
	TxSizeCostPerByte      uint64 `json:"tx_size_cost_per_byte" yaml:"tx_size_cost_per_byte"`
	SigVerifyCostED25519   uint64 `json:"sig_verify_cost_ed25519" yaml:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256k1 uint64 `json:"sig_verify_cost_secp256k1" yaml:"sig_verify_cost_secp256k1"`
}

// ------------------------------------------
//				Exported Account
// ------------------------------------------

type Account interface {
	GetAddress() sdk.AccAddress
	SetAddress(sdk.AccAddress) error // errors if already set.

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	GetCoins() sdk.Coins
	SetCoins(sdk.Coins) error

	// Calculates the amount of coins that can be sent to other accounts given
	// the current time.
	SpendableCoins(blockTime time.Time) sdk.Coins

	// Ensure that account implements stringer
	String() string
}

// GenesisAccounts defines a slice of GenesisAccount objects
type GenesisAccounts []GenesisAccount

// // Contains returns true if the given address exists in a slice of GenesisAccount
// // objects.
// func (ga GenesisAccounts) Contains(addr sdk.AccAddress) bool {
// 	for _, acc := range ga {
// 		if acc.GetAddress().Equals(addr) {
// 			return true
// 		}
// 	}

// 	return false
// }

// GenesisAccount defines a genesis account that embeds an Account with validation capabilities.
type GenesisAccount interface {
	Account
	Validate() error
}


// ------------------------------------------
//				BaseAccount
// ------------------------------------------


//-----------------------------------------------------------------------------
// BaseAccount

var _ Account = (*BaseAccount)(nil)
var _ GenesisAccount = (*BaseAccount)(nil)

// BaseAccount - a base account structure.
// This can be extended by embedding within in your AppAccount.
// However one doesn't have to use BaseAccount as long as your struct
// implements Account.
type BaseAccount struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        crypto.PubKey  `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
}

// // NewBaseAccount creates a new BaseAccount object
// func NewBaseAccount(address sdk.AccAddress, coins sdk.Coins,
// 	pubKey crypto.PubKey, accountNumber uint64, sequence uint64) *BaseAccount {

// 	return &BaseAccount{
// 		Address:       address,
// 		Coins:         coins,
// 		PubKey:        pubKey,
// 		AccountNumber: accountNumber,
// 		Sequence:      sequence,
// 	}
// }

// String implements fmt.Stringer
func (acc BaseAccount) String() string {
	var pubkey string

	if acc.PubKey != nil {
		pubkey = sdk.MustBech32ifyAccPub(acc.PubKey)
	}

// 	return fmt.Sprintf(`Account:
//   Address:       %s
//   Pubkey:        %s
//   Coins:         %s
//   AccountNumber: %d
//   Sequence:      %d`,
// 		acc.Address, pubkey, acc.Coins, acc.AccountNumber, acc.Sequence,
// 	)
// }

// // ProtoBaseAccount - a prototype function for BaseAccount
// func ProtoBaseAccount() exported.Account {
// 	return &BaseAccount{}
// }

// // NewBaseAccountWithAddress - returns a new base account with a given address
// func NewBaseAccountWithAddress(addr sdk.AccAddress) BaseAccount {
// 	return BaseAccount{
// 		Address: addr,
// 	}
// }

// GetAddress - Implements sdk.Account.
func (acc BaseAccount) GetAddress() sdk.AccAddress {
	return acc.Address
}

// SetAddress - Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr sdk.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// GetPubKey - Implements sdk.Account.
func (acc BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// SetPubKey - Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// GetCoins - Implements sdk.Account.
func (acc *BaseAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

// SetCoins - Implements sdk.Account.
func (acc *BaseAccount) SetCoins(coins sdk.Coins) error {
	acc.Coins = coins
	return nil
}

// GetAccountNumber - Implements Account
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// SpendableCoins returns the total set of spendable coins. For a base account,
// this is simply the base coins.
func (acc *BaseAccount) SpendableCoins(_ time.Time) sdk.Coins {
	return acc.GetCoins()
}

// // MarshalYAML returns the YAML representation of an account.
// func (acc BaseAccount) MarshalYAML() (interface{}, error) {
// 	var bs []byte
// 	var err error
// 	var pubkey string

// 	if acc.PubKey != nil {
// 		pubkey, err = sdk.Bech32ifyAccPub(acc.PubKey)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	bs, err = yaml.Marshal(struct {
// 		Address       sdk.AccAddress
// 		Coins         sdk.Coins
// 		PubKey        string
// 		AccountNumber uint64
// 		Sequence      uint64
// 	}{
// 		Address:       acc.Address,
// 		Coins:         acc.Coins,
// 		PubKey:        pubkey,
// 		AccountNumber: acc.AccountNumber,
// 		Sequence:      acc.Sequence,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return string(bs), err
// }

// Validate checks for errors on the account fields
func (acc BaseAccount) Validate() error {
	if acc.PubKey != nil && acc.Address != nil &&
		!bytes.Equal(acc.PubKey.Address().Bytes(), acc.Address.Bytes()) {
		return errors.New("pubkey and address pair is invalid")
	}

	return nil
}

// ------------------------------------------
//				Vesting Account
// ------------------------------------------

// ------------------------------------------
//				Codec
// ------------------------------------------

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*GenesisAccount)(nil), nil)
	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "cosmos-sdk/Account", nil)
	// cdc.RegisterConcrete(StdTx{}, "cosmos-sdk/StdTx", nil) // TODO ?
}