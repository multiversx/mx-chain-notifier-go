package notifier

import "errors"

// ErrNilTransactionPool signals that a nil transactions pool was provided
var ErrNilTransactionPool = errors.New("nil transactions pool")
