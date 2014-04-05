package libhydrogen

type Ledger struct {
	Accounts map[string]*Account
}

func NewLedger() *Ledger {
	return &Ledger{make(map[string]*Account)}
}
