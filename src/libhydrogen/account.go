package libhydrogen

type Account struct {
	ID       string
	Key      string
	Location string
	Balance  uint64
}

func (a *Account) Copy() *Account {
	return &Account{a.ID, a.Key, a.Location, a.Balance}
}
