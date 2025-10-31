package common

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/multiversx/mx-chain-core-go/data/stateChange"
)

// StateAccessToString returns state access as string
func StateAccessToString(stateAccess *stateChange.StateAccess) string {
	dataTrieChanges := make([]string, len(stateAccess.GetDataTrieChanges()))

	for i, dataTrieChange := range stateAccess.GetDataTrieChanges() {
		dataTrieChanges[i] = fmt.Sprintf("key: %v, val: %v, type: %v, operation %v, version %v", hex.EncodeToString(dataTrieChange.Key), hex.EncodeToString(dataTrieChange.Val), dataTrieChange.Type, dataTrieChange.Operation, dataTrieChange.Version)
	}

	return fmt.Sprintf("type: %v, operation: %v, mainTrieKey: %v, mainTrieVal: %v, index: %v, dataTrieChanges: %v, accountChanges %v",
		stateAccess.GetType(),
		stateAccess.GetOperation(),
		hex.EncodeToString(stateAccess.GetMainTrieKey()),
		hex.EncodeToString(stateAccess.GetMainTrieVal()),
		stateAccess.GetIndex(),
		strings.Join(dataTrieChanges, ", "),
		stateAccess.GetAccountChanges(),
	)
}
