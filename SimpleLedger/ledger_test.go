package SimpleLedger

import (
	"errors"
	"testing"
)

func TestLedgerTransaction(t *testing.T) {
	ledger := NewLedger()
	if err := ledger.AddAccount("abc", 100); err != nil {
		t.Error(err)
	}

	if err := ledger.AddAccount("abc", 100); !errors.Is(err, ErrAccountAlreadyExists) {
		t.Error(err)
	}

	if err := ledger.AddAccount("def", 100); err != nil {
		t.Error(err)
	}

	if err := ledger.DoTransaction("abc", "def", 50); err != nil {
		t.Error(err)
	}

	balance, err := ledger.GetBalance("abc")
	if err != nil {
		t.Error(err)
	}

	if balance != 50 {
		t.Error("Balance does not match")
	}

	balance, err = ledger.GetBalance("def")
	if err != nil {
		t.Error(err)
	}

	if balance != 150 {
		t.Error("Balance does not match")
	}
}
