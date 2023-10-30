package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	s := NewStore(testDb)

	a1, _, _ := createRandomAccount()
	a2, _, _ := createRandomAccount()
	fmt.Println(">>before: ", a1.Balance, a2.Balance)

	n := 5
	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := s.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: a1.ID,
				ToAccountID:   a2.ID,
				Amount:        amount,
			})

			errors <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, a1.ID, transfer.FromAccountID)
		require.Equal(t, a2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = s.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, a1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = s.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, a2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = s.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, a1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, a2.ID, toAccount.ID)

		fmt.Println(">>tx: ", fromAccount.Balance, toAccount.Balance)
		d1 := a1.Balance - fromAccount.Balance
		d2 := toAccount.Balance - a2.Balance
		require.Equal(t, d1, d2)
		require.True(t, d1 > 0)
		require.True(t, d1%amount == 0)

		k := int(d1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)

		existed[k] = true
	}

	ua1, err := testQueries.GetAccount(context.Background(), a1.ID)
	require.NoError(t, err)

	ua2, err := testQueries.GetAccount(context.Background(), a2.ID)
	require.NoError(t, err)

	fmt.Println(">>after: ", ua1.Balance, ua2.Balance)
	require.Equal(t, a1.Balance-int64(n)*amount, ua1.Balance)
	require.Equal(t, a2.Balance+int64(n)*amount, ua2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	s := NewStore(testDb)

	a1, _, _ := createRandomAccount()
	a2, _, _ := createRandomAccount()
	fmt.Println(">>before: ", a1.Balance, a2.Balance)

	n := 10
	amount := int64(10)

	errors := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := a1.ID
		toAccountId := a2.ID

		if i%2 == 1 {
			fromAccountID = a2.ID
			toAccountId = a1.ID
		}

		go func() {
			_, err := s.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			errors <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	ua1, err := testQueries.GetAccount(context.Background(), a1.ID)
	require.NoError(t, err)

	ua2, err := testQueries.GetAccount(context.Background(), a2.ID)
	require.NoError(t, err)

	fmt.Println(">>after: ", ua1.Balance, ua2.Balance)
	require.Equal(t, a1.Balance, ua1.Balance)
	require.Equal(t, a2.Balance, ua2.Balance)
}
