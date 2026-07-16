package telegram

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botmsg"
	"github.com/bots-go-framework/bots-fw/botplan"
	"github.com/bots-go-framework/bots-go-core/botkb"
)

// TelegramDescriptor is Telegram's static capability descriptor.
//
// Every field mirrors a record in can-i-use/capability-map.json (cited on the
// botplan.Descriptor field it sets). Telegram is the richest pilot platform:
// grid inline buttons, HTML with anchor links, edit and delete, callback acks,
// no send window.
var TelegramDescriptor = botplan.Descriptor{
	// Telegram's inline keyboard has no small hard button cap; a high value means
	// the prompt never degrades to a list (there is no list message on Telegram).
	MaxPromptButtons:    100, // telegram/inline-keyboard (grid, no documented small cap)
	MaxListRows:         0,   // telegram has no list message — not applicable
	MaxButtonLabelChars: 0,   // telegram/inline-keyboard: no documented label ceiling

	SupportsEdit:            true,  // telegram/edit-message native
	SupportsDelete:          true,  // telegram/delete-message native
	SupportsCallbackAck:     true,  // telegram/callback-query semanticsRequiresAck
	WindowGated:             false, // telegram/proactive-messaging native, no window
	SupportsInlineURLButton: true,  // telegram/inline-keyboard URL buttons
	SupportsButtonGrid:      true,  // telegram/inline-keyboard constraints.layout=grid
	SupportsMedia:           true,  // telegram/send-photo native
	TextMarkup:              botplan.MarkupHTML,
	SupportsAnchorTextLinks: true, // telegram/text-formatting hyperlinks=true
}

// Renderer renders a neutral botplan.MessagePlan into Telegram messages.
//
// It emits botmsg.MessageFromBot values the existing tgWebhookResponder already
// knows how to send: text uses FormatHTML (richToHTML), a prompt becomes an
// inline keyboard whose choice tokens ride as callback_data, a URL action
// becomes a URL button in the same keyboard, a live panel sets IsEdit, an image
// becomes a SendPhoto, and a proactive send is an ordinary message (Telegram has
// no window or templates).
type Renderer struct{}

var _ botplan.Renderer = Renderer{}

// NewRenderer returns a Telegram botplan.Renderer.
func NewRenderer() Renderer { return Renderer{} }

// Descriptor implements botplan.Renderer.
func (Renderer) Descriptor() botplan.Descriptor { return TelegramDescriptor }

// Render implements botplan.Renderer.
//
// Telegram carries a prompt and a URL action in a single keyboard, so a plan
// almost always renders to exactly one message. An image renders as a separate
// SendPhoto message preceding the text (Telegram photo captions are limited and
// the plan's text may be long), so a plan with Media yields two messages.
func (r Renderer) Render(plan botplan.MessagePlan, target botplan.RenderTarget) ([]botmsg.MessageFromBot, error) {
	if target.Platform != botplan.PlatformTelegram {
		return nil, botplan.ErrUnsupportedTarget
	}
	if err := plan.Validate(); err != nil {
		return nil, err
	}

	var msgs []botmsg.MessageFromBot

	// Image first, as its own message.
	if plan.Media != nil {
		msgs = append(msgs, r.photoMessage(plan.Media))
	}

	text := richToHTML(plan.Text)
	kb := r.keyboard(plan)

	// A media-only plan with empty text still needs no extra text message.
	if text == "" && kb == nil {
		return msgs, nil
	}

	m := botmsg.MessageFromBot{
		TextMessageFromBot: botmsg.TextMessageFromBot{
			Text:     text,
			Format:   botmsg.FormatHTML,
			Keyboard: kb,
			IsEdit:   plan.LivePanel != nil, // live-panel → edit-in-place
		},
	}
	msgs = append(msgs, m)
	return msgs, nil
}

// keyboard builds the inline keyboard for a plan's prompt and/or URL action, or
// nil if the plan has neither.
//
// Choices arrange into a grid honouring ActionPrompt.LayoutRows (a hint);
// choice tokens become callback_data via botkb.DataButton. A URL action is
// appended as a final row with a botkb.UrlButton.
func (Renderer) keyboard(plan botplan.MessagePlan) botkb.Keyboard {
	var rows [][]botkb.Button

	if plan.Prompt != nil {
		rows = append(rows, gridRows(plan.Prompt.Choices, plan.Prompt.LayoutRows)...)
	}
	if plan.URLAction != nil {
		rows = append(rows, []botkb.Button{botkb.NewUrlButton(plan.URLAction.Label, plan.URLAction.URL)})
	}
	if len(rows) == 0 {
		return nil
	}
	return botkb.NewMessageKeyboard(botkb.KeyboardTypeInline, rows...)
}

// gridRows arranges choices into rows honouring the layoutRows hint.
//
// layoutRows is the desired number of rows; buttons are spread across them as
// evenly as possible (a ceil-division per-row width). A zero or negative hint,
// or one exceeding the choice count, falls back to one button per row.
func gridRows(choices []botplan.Choice, layoutRows int) [][]botkb.Button {
	n := len(choices)
	if layoutRows <= 0 || layoutRows >= n {
		rows := make([][]botkb.Button, n)
		for i, c := range choices {
			rows[i] = []botkb.Button{botkb.NewDataButton(c.Label, c.Token)}
		}
		return rows
	}
	perRow := (n + layoutRows - 1) / layoutRows // ceil(n/layoutRows)
	var rows [][]botkb.Button
	for i := 0; i < n; i += perRow {
		end := i + perRow
		if end > n {
			end = n
		}
		row := make([]botkb.Button, 0, end-i)
		for _, c := range choices[i:end] {
			row = append(row, botkb.NewDataButton(c.Label, c.Token))
		}
		rows = append(rows, row)
	}
	return rows
}

// photoMessage builds a SendPhoto message from a MediaRef.
//
// MediaID is preferred (a Telegram file_id) when set; otherwise ImageURL is used
// as a PhotoUrl. The caption is plain text — the neutral Rich text stays on the
// following message.
func (Renderer) photoMessage(media *botplan.MediaRef) botmsg.MessageFromBot {
	var photo tgbotapi.Photo
	if media.MediaID != "" {
		photo = tgbotapi.FileID(media.MediaID)
	} else {
		photo = tgbotapi.PhotoUrl(media.ImageURL)
	}
	return botmsg.MessageFromBot{
		BotMessage: SendPhoto(tgbotapi.PhotoConfig{
			Photo:   photo,
			Caption: media.Caption,
		}),
	}
}
