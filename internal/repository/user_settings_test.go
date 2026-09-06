package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSettingsRepository_GetUserByToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserSettingsRepository(db)
	ctx := context.Background()

	t.Run("success finding user", func(t *testing.T) {
		futureDate := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
		rows := sqlmock.NewRows([]string{
			"id", "token", "password", "expired", "server_code", "quota", "relay", "adblock", "vpn",
		}).AddRow(1, "test-token", "pass123", futureDate, "SG", 10737418240, "", 1, "vmess")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, token, password, expired, server_code, quota, relay, adblock, vpn FROM users WHERE token = ? LIMIT 1;")).
			WithArgs("test-token").
			WillReturnRows(rows)

		u, err := repo.GetUserByToken(ctx, "test-token")
		require.NoError(t, err)
		assert.Equal(t, int64(1), u.ID)
		assert.Equal(t, "test-token", u.Token)
		assert.True(t, u.Adblock)
		assert.Equal(t, "vmess", u.VPN)
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, token, password, expired, server_code, quota, relay, adblock, vpn FROM users WHERE token = ? LIMIT 1;")).
			WithArgs("unknown-token").
			WillReturnError(sql.ErrNoRows)

		u, err := repo.GetUserByToken(ctx, "unknown-token")
		require.ErrorIs(t, err, repository.ErrUserNotFound)
		assert.Nil(t, u)
	})
}

func TestUserSettingsRepository_UpdateUserAdblock(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserSettingsRepository(db)
	ctx := context.Background()

	t.Run("success updating adblock to true", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET adblock = ? WHERE token = ?;")).
			WithArgs(true, "valid-token").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserAdblock(ctx, "valid-token", true)
		require.NoError(t, err)
	})

	t.Run("not found returns ErrUserNotFound", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET adblock = ? WHERE token = ?;")).
			WithArgs(false, "missing-token").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateUserAdblock(ctx, "missing-token", false)
		require.ErrorIs(t, err, repository.ErrUserNotFound)
	})
}

func TestUserSettingsRepository_UpdateUserUUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserSettingsRepository(db)
	ctx := context.Background()

	t.Run("success updating uuid", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET password = ? WHERE token = ?;")).
			WithArgs("new-uuid-1234", "my-token").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserUUID(ctx, "my-token", "new-uuid-1234")
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET password = ? WHERE token = ?;")).
			WithArgs("new-uuid-1234", "missing-token").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateUserUUID(ctx, "missing-token", "new-uuid-1234")
		require.ErrorIs(t, err, repository.ErrUserNotFound)
	})
}

func TestUserSettingsRepository_UpdateUserToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserSettingsRepository(db)
	ctx := context.Background()

	t.Run("success updating token", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET token = ? WHERE token = ?;")).
			WithArgs("new-tok", "old-tok").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserToken(ctx, "old-tok", "new-tok")
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET token = ? WHERE token = ?;")).
			WithArgs("new-tok", "old-tok").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateUserToken(ctx, "old-tok", "new-tok")
		require.ErrorIs(t, err, repository.ErrUserNotFound)
	})
}

func TestUserSettingsRepository_UpdateUserConfig(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserSettingsRepository(db)
	ctx := context.Background()

	t.Run("success updating vpn config", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET vpn = ?, server_code = ?, relay = ? WHERE token = ?;")).
			WithArgs("vless", "SG1", "SG", "my-token").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserConfig(ctx, "my-token", "vless", "SG1", "SG")
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET vpn = ?, server_code = ?, relay = ? WHERE token = ?;")).
			WithArgs("trojan", "ID1", "", "unknown-token").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateUserConfig(ctx, "unknown-token", "trojan", "ID1", "")
		require.ErrorIs(t, err, repository.ErrUserNotFound)
	})
}

func TestUserSettingsRepository_GetServers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserSettingsRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{
		"id", "code", "domain", "ip", "country", "users_count", "users_max",
	}).AddRow(1, "SG1", "sg1.example.com", "1.1.1.1", "SG", 10, 50).
		AddRow(2, "ID1", "id1.example.com", "2.2.2.2", "ID", 20, 100)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, code, domain, ip, country, users_count, users_max FROM servers ORDER BY code ASC;")).
		WillReturnRows(rows)

	servers, err := repo.GetServers(ctx)
	require.NoError(t, err)
	assert.Len(t, servers, 2)
	assert.Equal(t, "SG1", servers[0].Code)
	assert.Equal(t, "ID1", servers[1].Code)
}

func TestUserSettingsRepository_GetWildcards(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserSettingsRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"domain"}).
		AddRow("quiz.int.vidio.com").
		AddRow("zoom.us")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT domain FROM wildcards ORDER BY domain ASC;")).
		WillReturnRows(rows)

	wildcards, err := repo.GetWildcards(ctx)
	require.NoError(t, err)
	assert.Len(t, wildcards, 2)
	assert.Equal(t, "quiz.int.vidio.com", wildcards[0])
}

