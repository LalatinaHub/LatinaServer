package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
)

var (
	// ErrUserNotFound is returned when no user matches the given identifier.
	ErrUserNotFound = errors.New("user not found")
)

// UserSettingsRepository defines operations for managing user settings and portal features.
type UserSettingsRepository interface {
	GetUserByToken(ctx context.Context, token string) (*model.User, error)
	UpdateUserAdblock(ctx context.Context, token string, adblock bool) error
	UpdateUserUUID(ctx context.Context, token string, newUUID string) error
	UpdateUserToken(ctx context.Context, oldToken string, newToken string) error
	UpdateUserConfig(ctx context.Context, token string, vpn string, serverCode string, relay string) error
	GetServers(ctx context.Context) ([]model.Server, error)
	GetWildcards(ctx context.Context) ([]string, error)
}

type userSettingsRepo struct {
	db *sql.DB
}

// NewUserSettingsRepository creates a new UserSettingsRepository.
func NewUserSettingsRepository(db *sql.DB) UserSettingsRepository {
	return &userSettingsRepo{db: db}
}

func (r *userSettingsRepo) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	query := "SELECT id, token, password, expired, server_code, quota, relay, adblock, vpn FROM users WHERE token = ? LIMIT 1;"
	row := r.db.QueryRowContext(ctx, query, token)

	var (
		u          model.User
		expiredStr string
		adblockVal any
	)

	err := row.Scan(
		&u.ID,
		&u.Token,
		&u.Password,
		&expiredStr,
		&u.ServerCode,
		&u.Quota,
		&u.Relay,
		&adblockVal,
		&u.VPN,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user by token: %w", err)
	}

	if b, ok := adblockVal.(bool); ok {
		u.Adblock = b
	} else if i, ok := adblockVal.(int64); ok {
		u.Adblock = i > 0
	}

	parsedExpired, err := time.Parse("2006-01-02", expiredStr)
	if err == nil {
		u.Expired = parsedExpired
	} else if parsedRFC, err := time.Parse(time.RFC3339, expiredStr); err == nil {
		u.Expired = parsedRFC
	}

	return &u, nil
}

func (r *userSettingsRepo) UpdateUserAdblock(ctx context.Context, token string, adblock bool) error {
	query := "UPDATE users SET adblock = ? WHERE token = ?;"
	res, err := r.db.ExecContext(ctx, query, adblock, token)
	if err != nil {
		return fmt.Errorf("failed to update user adblock: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userSettingsRepo) UpdateUserUUID(ctx context.Context, token string, newUUID string) error {
	query := "UPDATE users SET password = ? WHERE token = ?;"
	res, err := r.db.ExecContext(ctx, query, newUUID, token)
	if err != nil {
		return fmt.Errorf("failed to update user uuid: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userSettingsRepo) UpdateUserToken(ctx context.Context, oldToken string, newToken string) error {
	query := "UPDATE users SET token = ? WHERE token = ?;"
	res, err := r.db.ExecContext(ctx, query, newToken, oldToken)
	if err != nil {
		return fmt.Errorf("failed to update user token: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userSettingsRepo) UpdateUserConfig(ctx context.Context, token string, vpn string, serverCode string, relay string) error {
	query := "UPDATE users SET vpn = ?, server_code = ?, relay = ? WHERE token = ?;"
	res, err := r.db.ExecContext(ctx, query, vpn, serverCode, relay, token)
	if err != nil {
		return fmt.Errorf("failed to update user config: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userSettingsRepo) GetServers(ctx context.Context) ([]model.Server, error) {
	query := "SELECT id, code, domain, ip, country, users_count, users_max FROM servers ORDER BY code ASC;"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query servers: %w", err)
	}
	defer rows.Close()

	var servers []model.Server
	for rows.Next() {
		var s model.Server
		if err := rows.Scan(&s.ID, &s.Code, &s.Domain, &s.IP, &s.Country, &s.UsersCount, &s.UsersMax); err != nil {
			return nil, fmt.Errorf("failed to scan server: %w", err)
		}
		servers = append(servers, s)
	}

	return servers, rows.Err()
}

func (r *userSettingsRepo) GetWildcards(ctx context.Context) ([]string, error) {
	query := "SELECT domain FROM wildcards ORDER BY domain ASC;"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		// If table does not exist or fails, return empty gracefully without breaking portal
		return []string{}, nil
	}
	defer rows.Close()

	var domains []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			continue
		}
		domains = append(domains, d)
	}

	return domains, nil
}

