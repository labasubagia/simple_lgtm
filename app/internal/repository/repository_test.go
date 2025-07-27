package repository_test

import (
	"context"
	"testing"

	"simple_lgtm/internal/model"
	"simple_lgtm/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository(t *testing.T) {

	t.Run("CreateData", func(t *testing.T) {
		ctx := context.Background()
		repo := repository.NewInMemoryRepository()

		err := repo.CreateData(ctx, "1", "value1")
		assert.NoError(t, err)

		err = repo.CreateData(ctx, "1", "value2")
		assert.Error(t, err)
	})

	t.Run("GetData", func(t *testing.T) {
		ctx := context.Background()
		repo := repository.NewInMemoryRepository()

		repo.CreateData(ctx, "2", "value2")
		value, err := repo.GetData(ctx, "2")
		assert.NoError(t, err)
		assert.Equal(t, "value2", value)

		_, err = repo.GetData(ctx, "nonexistent")
		assert.Error(t, err)
	})

	t.Run("UpdateData", func(t *testing.T) {
		ctx := context.Background()
		repo := repository.NewInMemoryRepository()

		repo.CreateData(ctx, "3", "value3")
		err := repo.UpdateData(ctx, "3", "newValue3")
		assert.NoError(t, err)

		value, _ := repo.GetData(ctx, "3")
		assert.Equal(t, "newValue3", value)

		err = repo.UpdateData(ctx, "nonexistent", "value")
		assert.Error(t, err)
	})

	t.Run("DeleteData", func(t *testing.T) {
		ctx := context.Background()
		repo := repository.NewInMemoryRepository()

		repo.CreateData(ctx, "4", "value4")
		err := repo.DeleteData(ctx, "4")
		assert.NoError(t, err)

		_, err = repo.GetData(ctx, "4")
		assert.Error(t, err)

		err = repo.DeleteData(ctx, "nonexistent")
		assert.Error(t, err)
	})

	t.Run("ListAllData", func(t *testing.T) {
		ctx := context.Background()
		repo := repository.NewInMemoryRepository()

		repo.CreateData(ctx, "5", "value5")
		repo.CreateData(ctx, "6", "value6")

		dataList, err := repo.ListAllData(ctx)
		assert.NoError(t, err)
		assert.Len(t, dataList, 2)

		assert.Contains(t, dataList, model.DataItem{ID: "5", Value: "value5"})
		assert.Contains(t, dataList, model.DataItem{ID: "6", Value: "value6"})
	})
}
