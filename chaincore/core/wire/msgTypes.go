package wire

import "strconv"

const (
	VERSION        = "version"
	VERACK         = "verack"
	ADDR           = "addr"
	ADDRV2         = "addrv2"
	SENDADDRV2     = "sendaddrv2"
	INV            = "inv"
	GETDATA        = "getdata"
	MERKLEBLOCK    = "merkleblock"
	GETBLOCKS      = "getblocks"
	GETHEADERS     = "getheaders"
	TX             = "tx"
	HEADERS        = "headers"
	BLOCK          = "block"
	GETADDR        = "getaddr"
	MEMPOOL        = "mempool"
	PING           = "ping"
	PONG           = "pong"
	NOTFOUND       = "notfound"
	FILTERLOAD     = "filterload"
	FILTERADD      = "filteradd"
	FILTERCLEAR    = "filterclear"
	SENDHEADERS    = "sendheaders"
	FEEFILTER      = "feefilter"
	SENDCMPCT      = "sendcmpct"
	CMPCTBLOCK     = "cmpctblock"
	GETBLOCKTXN    = "getblocktxn"
	BLOCKTXN       = "blocktxn"
	GETCFILTERS    = "getcfilters"
	CFILTER        = "cfilter"
	GETCFHEADERS   = "getcfheaders"
	CFHEADERS      = "cfheaders"
	GETCFCHEKCPT   = "getcfcheckpt"
	CFCHECKPT      = "cfcheckpt"
	ADDVALIDATOR   = "addvalidator"
	GETVALIDATORS  = "getvalidator"
	STAKE          = "stake"
	UNSTAKE        = "unstake"
	ANONPROXYRELAY = "anonproxyrelay"
)

const (
	SERVICE_NODE_NONE            = 0
	SERVICE_NODE_NETWORK         = 1 << 0
	SERVICE_NODE_UTXO            = 1 << 1
	SERVICE_NODE_BLOOM           = 1 << 2
	SERVICE_NODE_VALIDATOR       = 1 << 3
	SERVICE_NODE_NETWORK_LIMITED = 1 << 10
	SERVICE_NODE_WVM             = 1 << 6
)

func ServiceFlagsToString(sflags ...uint64) []string {
	var strs []string
	strs = make([]string, len(sflags))
	for _, flag := range sflags {

		strs = append(strs, strconv.FormatUint(flag, 10))
	}

	return strs
}
