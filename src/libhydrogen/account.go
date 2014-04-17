package libhydrogen

type Account struct {
	Key      string
	Location string
	Balance  uint64
}

func (a *Account) Copy() *Account {
	return &Account{a.Key, a.Location, a.Balance}
}
