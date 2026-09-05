package db

import (
	"context"
	"strconv"

	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/LalatinaHub/LatinaServer/internal/database"
	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/LalatinaHub/LatinaServer/internal/repository"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// Legacy type aliases for backward compatibility
type UserStruct = model.User
type ServerStruct = model.Server

func GetKVList() map[string]any {
	ctx := context.Background()
	db, err := database.GetDB()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get database connection")
		return make(map[string]any)
	}
	kvRepo := repository.NewKVRepository(db)

	kvList, err := kvRepo.GetAll(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get KV list")
		return make(map[string]any)
	}

	return kvList
}

func GetServerList() []model.Server {
	ctx := context.Background()
	db, err := database.GetDB()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get database connection")
		return nil
	}
	serverRepo := repository.NewServerRepository(db)

	serverList, err := serverRepo.GetAll(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get server list")
		return nil
	}

	return serverList
}

func GetPremiumList() map[string][]model.User {
	ctx := context.Background()
	db, err := database.GetDB()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get database connection")
		return make(map[string][]model.User)
	}
	userRepo := repository.NewUserRepository(db)

	userMap, err := userRepo.GetActiveUsersGroupedByVPN(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get premium user list")
		return make(map[string][]model.User)
	}

	return userMap
}

func UpdateAndCheckPremiumQuota(name string) bool {
	ctx := context.Background()
	db, err := database.GetDB()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get database connection for quota check")
		return true
	}
	userRepo := repository.NewUserRepository(db)

	id, err := strconv.ParseInt(name, 10, 64)
	if err != nil {
		logger.Error().Err(err).Str("name", name).Msg("Failed to parse user ID")
		return true
	}

	usedBytes := helper.GetUserStats(name)
	_, quotaExhausted, err := userRepo.DeductQuota(ctx, id, usedBytes)
	if err != nil {
		logger.Error().Err(err).Int64("user_id", id).Msg("Failed to deduct quota")
		return true // Treat error as quota exhausted for safety
	}

	return quotaExhausted
}


