package telegram

import (
	"bytes"
	"context"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/bots-go-framework/bots-fw/mocks/mock_botsfw"
	"github.com/strongo/i18n"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

func TestNewTelegramWebhookHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewTelegramWebhookHandler() did not panic")
		}
	}()
	NewTelegramWebhookHandler(nil, nil, nil)
}

func TestTelegramWebhookHandler_Handle(t *testing.T) {
	t.Run("SuccessfulPayment", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		botContextProvider := mock_botsfw.NewMockBotContextProvider(mockCtrl)

		ctx := context.Background()

		botContext := botsfw.BotContext{}
		botContextProvider.EXPECT().GetBotContext(ctx, PlatformID, gomock.Any()).Return(&botContext, nil)

		var translatorProvider botsfw.TranslatorProvider = func(c context.Context) i18n.Translator {
			return nil
		}
		setAppUserFields := func(botsfwmodels.AppUserData, botinput.Sender) error {
			return nil
		}
		handler := NewTelegramWebhookHandler(botContextProvider, translatorProvider, setAppUserFields)
		var r http.Request
		r.Method = "POST"

		// Create the HTTP request with context
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://example.com/api", bytes.NewReader([]byte(`
{
    "update_id": 903594254,
    "message": {
        "message_id": 3648,
        "from": {
            "id": 92819884,
            "is_bot": false,
            "first_name": "Alexander",
            "last_name": "Trakhimenok",
            "username": "trakhimenok",
            "language_code": "en"
        },
        "chat": {
            "id": 92819884,
            "first_name": "Alexander",
            "last_name": "Trakhimenok",
            "username": "trakhimenok",
            "type": "private"
        },
        "date": 1750355864,
        "successful_payment": {
            "currency": "XTR",
            "total_amount": 2,
            "invoice_payload": "topped_up",
            "telegram_payment_charge_id": "some_charge_id",
            "provider_payment_charge_id": "1234567890_12"
        }
    }
}
`)))

		if err != nil {
			t.Fatalf("NewRequestWithContext failed: %v", err)
		}

		_, entriesWithInputs, err := handler.GetBotContextAndInputs(ctx, req)
		if err != nil {
			t.Fatalf("handler.GetBotContextAndInputs failed: %v", err)
		}
		if len(entriesWithInputs) != 1 {
			t.Errorf("len(entriesWithInputs) = %v, want 1", len(entriesWithInputs))
		}
		if len(entriesWithInputs[0].Inputs) != 1 {
			t.Errorf("len(entriesWithInputs[0].Inputs) = %v, want 1", len(entriesWithInputs[0].Inputs))
		}
		if inputType := entriesWithInputs[0].Inputs[0].InputType(); inputType != botinput.TypeSuccessfulPayment {
			t.Errorf("entriesWithInputs[0].Inputs[0].InputType() = %v, want %v", inputType, botinput.TypeSuccessfulPayment)
		}
		successfulPayment := entriesWithInputs[0].Inputs[0].(botinput.SuccessfulPayment)

		if currency := successfulPayment.GetCurrency(); currency != "XTR" {
			t.Errorf("successfulPayment.GetCurrency() = %v, want XTR", currency)
		}
		if totalAmount := successfulPayment.GetTotalAmount(); totalAmount != 2 {
			t.Errorf("successfulPayment.GetTotalAmount() = %v, want 2", totalAmount)
		}
		if messengerChargeID := successfulPayment.GetMessengerChargeID(); messengerChargeID != "some_charge_id" {
			t.Errorf("successfulPayment.GetMessengerChargeID() = %v, want some_charge_id", messengerChargeID)
		}
		if paymentProviderChargeID := successfulPayment.GetPaymentProviderChargeID(); paymentProviderChargeID != "1234567890_12" {
			t.Errorf("successfulPayment.GetPaymentProviderChargeID() = %v, want some_charge_id", "1234567890_12")
		}
	})
}
