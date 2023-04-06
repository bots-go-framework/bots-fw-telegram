package telegram

import (
	"github.com/bots-go-framework/bots-fw-dalgo/dalgo4botsfw"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDalgoStores(t *testing.T) {
	type args struct {
		db dalgo4botsfw.DbProvider
	}
	tests := []struct {
		name  string
		args  args
		want  botsfw.BotChatStore
		want1 botsfw.BotUserStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewDalgoStores(tt.args.db)
			assert.Equalf(t, tt.want, got, "NewDalgoStores(%v)", tt.args.db)
			assert.Equalf(t, tt.want1, got1, "NewDalgoStores(%v)", tt.args.db)
		})
	}
}
