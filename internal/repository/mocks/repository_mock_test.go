package mocks_test

import (
	"context"
	"testing"

	"github.com/LalatinaHub/LatinaServer/internal/domain/model"
	"github.com/LalatinaHub/LatinaServer/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMockUserRepository(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	ctx := context.Background()

	expectedUsers := map[string][]model.User{
		"vmess": {{ID: 1, Token: "t1"}},
	}

	mockRepo.On("GetActiveUsersGroupedByVPN", ctx).Return(expectedUsers, nil)
	mockRepo.On("DeductQuota", ctx, int64(1), int64(100)).Return(int64(900), false, nil)
	mockRepo.On("DeductQuotaBatch", ctx, map[int64]int64{1: 100}).Return([]int64{}, nil)

	users, err := mockRepo.GetActiveUsersGroupedByVPN(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)

	remaining, depleted, err := mockRepo.DeductQuota(ctx, 1, 100)
	assert.NoError(t, err)
	assert.Equal(t, int64(900), remaining)
	assert.False(t, depleted)

	depletedUsers, err := mockRepo.DeductQuotaBatch(ctx, map[int64]int64{1: 100})
	assert.NoError(t, err)
	assert.Empty(t, depletedUsers)

	mockRepo.AssertExpectations(t)
}
