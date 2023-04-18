package telegram

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	t.Run("panics if getDb == nil", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("did not panic")
			}
		}()
		Init(nil)
	})
	t.Run("sets getDatabase", func(t *testing.T) {
		getDb := func(context.Context) (dal.Database, error) {
			return nil, nil
		}
		assert.Nil(t, getDatabase)
		Init(getDb)
		assert.NotNil(t, getDatabase)
	})
}
