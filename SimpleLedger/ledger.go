package SimpleLedger

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

var (
	ErrAccountAlreadyExists         = errors.New("account already exists")
	ErrAccountDoesNotExist          = errors.New("account does not exist")
	ErrAmountIsNotEnoughForTransfer = errors.New("amount is not enough for transfer")
)

type Entry struct {
	accountID string
	amount    float64
}
type Transaction struct {
	id      string
	entries []Entry
}

type Ledger struct {
	mu                        sync.Mutex
	listOfAccountsWithBalance map[string]float64
	listOfTransactions        []Transaction
}

func NewLedger() *Ledger {
	return &Ledger{
		listOfAccountsWithBalance: make(map[string]float64),
		listOfTransactions:        make([]Transaction, 0),
	}
}

func (ledger *Ledger) AddAccount(accountID string, amount float64) error {
	_, isExist := ledger.listOfAccountsWithBalance[accountID]
	if isExist {
		return ErrAccountAlreadyExists
	}

	ledger.mu.Lock()
	defer ledger.mu.Unlock()

	ledger.listOfAccountsWithBalance[accountID] = amount
	return nil
}

func (ledger *Ledger) RemoveAccount(accountID string) error {
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	_, isExist := ledger.listOfAccountsWithBalance[accountID]
	if isExist {
		delete(ledger.listOfAccountsWithBalance, accountID)
		return nil
	}
	return ErrAccountDoesNotExist
}

func (ledger *Ledger) GetBalance(accountID string) (float64, error) {
	ledger.mu.Lock()
	defer ledger.mu.Unlock()

	if amount, isExist := ledger.listOfAccountsWithBalance[accountID]; isExist {
		return amount, nil
	}

	return 0, ErrAccountDoesNotExist
}

func (ledger *Ledger) DoTransaction(from, to string, amount float64) error {
	ledger.mu.Lock()
	defer ledger.mu.Unlock()

	if ledger.listOfAccountsWithBalance[from] < amount {
		return ErrAmountIsNotEnoughForTransfer
	}

	tx := Transaction{
		id: uuid.New().String(),
		entries: []Entry{
			{
				accountID: from,
				amount:    -amount,
			},
			{
				accountID: to,
				amount:    amount,
			},
		},
	}

	for _, entry := range tx.entries {
		ledger.listOfAccountsWithBalance[entry.accountID] += entry.amount
	}

	ledger.listOfTransactions = append(ledger.listOfTransactions, tx)
	return nil
}
