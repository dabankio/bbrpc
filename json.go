package bbrpc

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// Request is a type for raw JSON-RPC 1.0 requests.  The Method field identifies
// the specific command type which in turns leads to different parameters.
// Callers typically will not use this directly since this package provides a
// statically typed command infrastructure which handles creation of these
// requests, however this struct it being exported in case the caller wants to
// construct raw requests for some reason.
type Request struct {
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

// Keypair .
type Keypair struct {
	Privkey string
	Pubkey  string
}

// AddrKeypair address and keypair
type AddrKeypair struct {
	Keypair
	Address string
}

// BalanceInfo .
type BalanceInfo struct {
	Address     string  `json:"address,omitempty"`     //(string, required) wallet address
	Avail       float64 `json:"avail,omitempty"`       //(double, required) balance available amount
	Locked      float64 `json:"locked,omitempty"`      //(double, required) locked amount
	Unconfirmed float64 `json:"unconfirmed,omitempty"` //(double, required) unconfirmed amount
}

// PubkeyInfo .
type PubkeyInfo struct {
	Key     string `json:"key,omitempty"` //public key with hex system
	Version int64  `json:"version,omitempty"`
	Locked  bool   `json:"locked,omitempty"`
	Timeout *int64 `json:"timeout,omitempty"` //public key timeout locked
}

// TemplateData .
type TemplateData struct {
	Type string `json:"type,omitempty"` //(string, required) template type
	Hex  string `json:"hex,omitempty"`  //(string, required) temtplate data
	//   (if type=delegate)
	Delegate *struct { //(object, required) delegate template struct
		Delegate string `json:"delegate,omitempty"` //(string, required) delegate public key
		Owner    string `json:"owner,omitempty"`    //(string, required) owner address
	} `json:"delegate,omitempty"`
	//   (if type=vote)
	Vote *struct { //(object, required) delegate template struct
		Delegate string `json:"delegate,omitempty"` //(string, required) delegate public key
		Owner    string `json:"owner,omitempty"`    //(string, required) owner address
	} `json:"vote,omitempty"`
	//   (if type=fork)
	Fork *struct { //(object, required) fork template struct
		Redeem string `json:"redeem,omitempty"` //(string, required) redeem address
		Fork   string `json:"fork,omitempty"`   //(string, required) fork hash
	} `json:"fork,omitempty"`
	//   (if type=mint)
	Mint *struct { //(object, required) mint template struct
		Mint  string `json:"mint,omitempty"`  //(string, required) mint public key
		Spent string `json:"spent,omitempty"` //(string, required) spent address
	} `json:"mint,omitempty"`
	//   (if type=multisig)
	Multisig *struct { //(object, required) multisig template struct
		Required  int64    `json:"required,omitempty"`  //(int, required) required weight
		Addresses []string `json:"addresses,omitempty"` //(string, required) public key
	} `json:"multisig,omitempty"`
	//   (if type=exchange)
	Exchange *struct { //(object, required) exchange template struct
		SpendS  string `json:"spend_s,omitempty"`  //(string, required) spend_s
		SpendM  string `json:"spend_m,omitempty"`  //(string, required) spend_m
		HeightM int64  `json:"height_m,omitempty"` //(int, required) height m
		HeightS int64  `json:"height_s,omitempty"` //(int, required) height s
		ForkM   string `json:"fork_m,omitempty"`   //(string, required) fork m
		ForkS   string `json:"fork_s,omitempty"`   //(string, required) fork s
	} `json:"exchange,omitempty"`
	//   (if type=weighted)
	Weighted *struct { //(object, required) weighted template struct
		Required int64      `json:"required,omitempty"` //(int, required) required weight
		Pubkey   []struct { //(object, required) public key
			Key    string `json:"key,omitempty"`    //(string, required) public key
			Weight int64  `json:"weight,omitempty"` //(int, required) weight
		} `json:"pubkey,omitempty"`
	} `json:"weighted,omitempty"`
	Dexbbcmap *struct {
		Owner string `json:"owner"`
	} `json:"dexbbcmap,omitempty"`
}

// AddressData .
type AddressData struct {
	Address string `json:"address,omitempty"` //(string, required) wallet address
	Ismine  bool   `json:"ismine,omitempty"`  //(bool, required) is mine
	Type    string `json:"type,omitempty"`    //(string, required) type, pubkey or template
	// (if type=pubkey)
	Pubkey *string `json:"pubkey,omitempty"` //(string, required) public key
	// (if type=template)
	Template *string `json:"template,omitempty"` //(string, required) template type name
	// (if type=template && ismine=true)
	Templatedata *TemplateData `json:"templatedata,omitempty"` //(object, required) template data
}

// AddressInfo .
type AddressInfo struct {
	Isvalid bool `json:"isvalid,omitempty"` //": true|false,               (bool, required) is valid
	//    (if isvalid=true)
	Addressdata *AddressData `json:"addressdata,omitempty"` //(object, required) address data
}

// AddnewtemplateParam .
type AddnewtemplateParam interface {
	TemplateType() string
	ParamName() string
}

// AddnewtemplateParamMint .
type AddnewtemplateParamMint struct {
	Mint  string `json:"mint,omitempty"`  //mint public key
	Spent string `json:"spent,omitempty"` //spent address
}

// TemplateType .
func (p AddnewtemplateParamMint) TemplateType() string { return "mint" }

// ParamName .
func (p AddnewtemplateParamMint) ParamName() string { return "mint" }

// AddnewtemplateParamDexBBCMap .
type AddnewtemplateParamDexBBCMap struct {
	Owner    string `json:"owner,omitempty"`    //(string, required) owner address
}

// TemplateType .
func (p AddnewtemplateParamDexBBCMap) TemplateType() string { return "dexbbcmap" }

// ParamName .
func (p AddnewtemplateParamDexBBCMap) ParamName() string { return "dexbbcmap" }

// AddnewtemplateParamDelegate .
type AddnewtemplateParamDelegate struct {
	Delegate string `json:"delegate,omitempty"` //(string, required) delegate public key
	Owner    string `json:"owner,omitempty"`    //(string, required) owner address
}

// TemplateType .
func (p AddnewtemplateParamDelegate) TemplateType() string { return "delegate" }

// ParamName .
func (p AddnewtemplateParamDelegate) ParamName() string { return "delegate" }

// AddnewtemplateParamVote .
type AddnewtemplateParamVote struct {
	Delegate string `json:"delegate,omitempty"` //(string, required) delegate public key
	Owner    string `json:"owner,omitempty"`    //(string, required) owner address
}

// TemplateType .
func (p AddnewtemplateParamVote) TemplateType() string { return "vote" }

// ParamName .
func (p AddnewtemplateParamVote) ParamName() string { return "vote" }

// AddnewtemplateParamMultisig .
type AddnewtemplateParamMultisig struct {
	Required int      `json:"required,omitempty"`
	Pubkeys  []string `json:"pubkeys,omitempty"`
}

// TemplateType .
func (p AddnewtemplateParamMultisig) TemplateType() string { return "multisig" }

// ParamName .
func (p AddnewtemplateParamMultisig) ParamName() string { return "multisig" }

// ForkProfile .
type ForkProfile struct {
	Fork       string  `json:"fork,omitempty"`       //(string, required) fork id with hex system
	Name       string  `json:"name,omitempty"`       //(string, required) fork name
	Symbol     string  `json:"symbol,omitempty"`     //(string, required) fork symbol
	Amount     float64 `json:"amount,omitempty"`     //(double, required) amount
	Reward     float64 `json:"reward,omitempty"`     //(double, required) mint reward
	Halvecycle uint64  `json:"halvecycle,omitempty"` //(uint, required) halve cycle: 0: fixed reward, >0: blocks of halve cycle
	Isolated   bool    `json:"isolated,omitempty"`   //(bool, required) is isolated
	Private    bool    `json:"private,omitempty"`    //(bool, required) is private
	Enclosed   bool    `json:"enclosed,omitempty"`   //(bool, required) is enclosed
	Owner      string  `json:"owner,omitempty"`      //(string, required) owner's address
}

// Delegate .
type Delegate struct {
	Address string  `json:"address,omitempty"`
	Votes   float64 `json:"votes,omitempty"`
}

// BlockInfo 块信息
type BlockInfo struct {
	Hash    string   `json:"hash"`    //(string, required) block hash
	Version uint     `json:"version"` //(uint, required) version
	Type    string   `json:"type"`    //(string, required) block type
	Time    uint     `json:"time"`    //(uint, required) block time
	Fork    string   `json:"fork"`    //(string, required) fork hash
	Height  uint     `json:"height"`  //(uint, required) block height
	Txmint  string   `json:"txmint"`  //(string, required) transaction mint hash
	Prev    string   `json:"prev"`    //(string, optional) previous block hash
	Tx      []string `json:"tx"`      //(string, required) transaction hash
}

// BlockDetail 区块详细信息
type BlockDetail struct {
	Hash     string                      `json:"hash,omitempty"`     //(string, required) block hash
	HashPrev string                      `json:"hashPrev,omitempty"` //(string, required) block prev hash
	Version  uint                        `json:"version,omitempty"`  //(uint, required) version
	Type     string                      `json:"type,omitempty"`     //(string, required) block type
	Time     uint                        `json:"time,omitempty"`     //(uint, required) block time
	Bits     uint                        `json:"bits,omitempty"`     //(uint, required) nBits
	Fork     string                      `json:"fork,omitempty"`     //(string, required) fork hash
	Height   uint                        `json:"height,omitempty"`   //(uint, required) block height
	Prev     string                      `json:"prev,omitempty"`     //(string, optional) previous block hash
	Txmint   NoneSerializedTransaction   `json:"txmint,omitempty"`   // (object, required) transaction mint data
	Tx       []NoneSerializedTransaction `json:"tx,omitempty"`
}

// CmdCreatetransaction .
type CmdCreatetransaction struct {
	From   string   `json:"from,omitempty"`   //(string, required) from address
	To     string   `json:"to,omitempty"`     //(string, required) to address
	Amount float64  `json:"amount,omitempty"` //(double, required) amount
	Txfee  *float64 `json:"txfee,omitempty"`  //(double, optional) transaction fee
	Fork   *string  `json:"fork,omitempty"`   //(string, optional) fork hash
	Data   *string  `json:"data,omitempty"`   //(string, optional) output data
}

// SigntransactionResult .
type SigntransactionResult struct {
	Hex       string `json:"hex,omitempty"`       //(string, required) hex of transaction data
	Completed bool   `json:"completed,omitempty"` //(bool, required) transaction completed or not
}

// Transaction .
type Transaction struct {
	Txid        string  `json:"txid,omitempty"`        //(string, required) transaction hash
	Fork        string  `json:"fork,omitempty"`        //(string, required) fork hash
	Type        string  `json:"type,omitempty"`        //(string, required) transaction type
	Time        uint    `json:"time,omitempty"`        //(uint, required) transaction timestamp
	Send        bool    `json:"send,omitempty"`        //(bool, required) is from me
	To          string  `json:"to,omitempty"`          //(string, required) to address
	Amount      float64 `json:"amount,omitempty"`      //(double, required) transaction amount
	Fee         float64 `json:"fee,omitempty"`         //(double, required) transaction fee
	Lockuntil   uint    `json:"lockuntil,omitempty"`   //(uint, required) lockuntil
	Blockheight *int    `json:"blockheight,omitempty"` //(int, optional) block height
	From        *string `json:"from,omitempty"`        //(string, optional) from address
}

// TransactionDetail .
type TransactionDetail struct {
	// (if serialized=true)
	Serialization *string `json:"serialization,omitempty"` //(string, optional) transaction hex data
	//    (if serialized=false)
	Transaction *NoneSerializedTransaction `json:"transaction,omitempty"` //(object, optional) transaction data
}

// NoneSerializedTransaction .
type NoneSerializedTransaction struct {
	Txid          string     `json:"txid,omitempty"`          //(string, required) transaction hash
	Version       uint       `json:"version,omitempty"`       //(uint, required) version
	Type          string     `json:"type,omitempty"`          //(string, required) transaction type
	Time          uint       `json:"time,omitempty"`          //(uint, required) transaction timestamp
	Lockuntil     uint       `json:"lockuntil,omitempty"`     //(uint, required) unlock time
	Anchor        string     `json:"anchor,omitempty"`        //(string, required) anchor hash
	Blockhash     string     `json:"blockhash"`               //block hash
	Sendfrom      string     `json:"sendfrom,omitempty"`      //(string, required) send from address
	Sendto        string     `json:"sendto,omitempty"`        //(string, required) send to address
	Amount        float64    `json:"amount,omitempty"`        //(double, required) amount
	Txfee         float64    `json:"txfee,omitempty"`         //(double, required) transaction fee
	Data          string     `json:"data,omitempty"`          //(string, required) data
	Sig           string     `json:"sig,omitempty"`           //(string, required) sign
	Fork          string     `json:"fork,omitempty"`          //(string, required) fork hash
	Confirmations *int       `json:"confirmations,omitempty"` //(int, optional) confirmations
	Vin           []VinPoint `json:"vin,omitempty"`           //(object, required) vin struct
}

// VinPoint .
type VinPoint struct {
	Txid string `json:"txid,omitempty"` //(string, required) pre-vout transaction hash
	Vout uint   `json:"vout,omitempty"` //(uint, required) pre-vout number
}

// PeerInfo .
type PeerInfo struct {
	Address  string `json:"address,omitempty"`  //(string, required) peer address
	Services string `json:"services,omitempty"` //(string, required) service
	Lastsend string `json:"lastsend,omitempty"` //(int, required) last send time
	Lastrecv string `json:"lastrecv,omitempty"` //(int, required) last receive time
	Conntime string `json:"conntime,omitempty"` //(int, required) active time
	Pingtime int    `json:"pingtime"`           //(int, required) ping pong time
	Version  string `json:"version,omitempty"`  //(string, required) version
	Subver   string `json:"subver,omitempty"`   //(string, required) sub version
	Inbound  bool   `json:"inbound,omitempty"`  //(bool, required) accept multiple connection or not
	Height   int    `json:"height,omitempty"`   //(int, required) starting height
	Banscore int    `json:"banscore,omitempty"` //(int, required) ban score
}

// UnspentTotal .
type UnspentTotal struct {
	Total     decimal.Decimal `json:"total,omitempty"`
	Addresses []struct {
		Address  string `json:"address,omitempty"`
		Unspents []struct {
			Txid      string          `json:"txid,omitempty"`
			Out       uint            `json:"out,omitempty"`
			Amount    decimal.Decimal `json:"amount,omitempty"`
			Time      int64           `json:"time,omitempty"`
			Lockuntil int64           `json:"lockuntil,omitempty"`
		} `json:"unspents,omitempty"`
		Sum decimal.Decimal `json:"sum,omitempty"`
	} `json:"addresses,omitempty"`
}
