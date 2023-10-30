package db

import (
	"context"
	"database/sql"
	"github.com/VrelinnVailati/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount() (Account, CreateAccountParams, error) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)

	return account, arg, err
}

func TestCreateAccount(t *testing.T) {
	a, arg, err := createRandomAccount()

	require.NoError(t, err)
	require.NotEmpty(t, a)

	require.Equal(t, arg.Owner, a.Owner)
	require.Equal(t, arg.Balance, a.Balance)
	require.Equal(t, arg.Currency, a.Currency)

	require.NotZero(t, a.ID)
	require.NotZero(t, a.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	a, _, _ := createRandomAccount()
	ga, err := testQueries.GetAccount(context.Background(), a.ID)

	require.NoError(t, err)
	require.NotEmpty(t, ga)

	require.Equal(t, a.ID, ga.ID)
	require.Equal(t, a.Owner, ga.Owner)
	require.Equal(t, a.Currency, ga.Currency)
	require.Equal(t, a.Balance, ga.Balance)
	require.WithinDuration(t, a.CreatedAt, ga.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	a, _, _ := createRandomAccount()

	uarg := UpdateAccountParams{
		ID:      a.ID,
		Balance: util.RandomMoney(),
	}

	ua, err := testQueries.UpdateAccount(context.Background(), uarg)
	require.NoError(t, err)
	require.NotEmpty(t, ua)

	require.Equal(t, a.ID, ua.ID)
	require.Equal(t, a.Owner, ua.Owner)
	require.Equal(t, a.Currency, ua.Currency)
	require.Equal(t, uarg.Balance, ua.Balance)
	require.WithinDuration(t, a.CreatedAt, ua.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	a, _, _ := createRandomAccount()

	err := testQueries.DeleteAccount(context.Background(), a.ID)
	require.NoError(t, err)

	da, err := testQueries.GetAccount(context.Background(), a.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, da)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount()
	}

	largs := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), largs)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, a := range accounts {
		require.NotEmpty(t, a)
	}
}
