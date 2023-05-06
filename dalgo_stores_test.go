package telegram

import (
	"context"
	"github.com/bots-go-framework/dalgo4botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDalgoStores(t *testing.T) {
	type args struct {
		db dalgo4botsfw.DbProvider
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "empty", args: args{db: nil}},
		{name: "should_pass", args: args{db: func(c context.Context) (dal.Database, error) {
			return nil, nil
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.db == nil {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("The code did not panic")
					}
				}()
			}
			chatStore, botUserStore := NewDalgoStores(tt.args.db)
			assert.NotNil(t, chatStore, "chatStore")
			assert.NotNil(t, botUserStore, "botUserStore")
		})
	}
}
