package telegram

import (
	"errors"
	"testing"

	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botmsg"
	"github.com/bots-go-framework/bots-fw/botplan"
	"github.com/bots-go-framework/bots-go-core/botkb"
)

var tgTarget = botplan.RenderTarget{Platform: botplan.PlatformTelegram, Locale: "en"}

func TestRichToHTML(t *testing.T) {
	tests := []struct {
		name string
		rich botplan.Rich
		want string
	}{
		{
			name: "bold and italic",
			rich: botplan.Rich{Lines: []botplan.Line{botplan.Para(
				botplan.Text("go "), botplan.Bold("fast"), botplan.Text(" & "), botplan.Italic("slow"))}},
			want: "go <b>fast</b> &amp; <i>slow</i>",
		},
		{
			name: "escaping of literal angle brackets",
			rich: botplan.RichText("a < b > c & d"),
			want: "a &lt; b &gt; c &amp; d",
		},
		{
			name: "code span",
			rich: botplan.Rich{Lines: []botplan.Line{botplan.Para(botplan.Code("go test"))}},
			want: "<code>go test</code>",
		},
		{
			name: "anchor link keeps anchor text",
			rich: botplan.Rich{Lines: []botplan.Line{botplan.Para(botplan.Link("View day", "https://x.io/d?a=1&b=2"))}},
			want: `<a href="https://x.io/d?a=1&amp;b=2">View day</a>`,
		},
		{
			name: "list items and quote",
			rich: botplan.Rich{Lines: []botplan.Line{
				botplan.Item(botplan.Text("kite")),
				botplan.Quote(botplan.Text("moved")),
			}},
			want: "• kite\n<blockquote>moved</blockquote>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := richToHTML(tt.rich); got != tt.want {
				t.Errorf("richToHTML() =\n  %q\nwant\n  %q", got, tt.want)
			}
		})
	}
}

func TestRenderTextOnly(t *testing.T) {
	plan := botplan.MessagePlan{Text: botplan.RichText("hello")}
	msgs, err := NewRenderer().Render(plan, tgTarget)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 {
		t.Fatalf("want 1 message, got %d", len(msgs))
	}
	m := msgs[0]
	if m.Text != "hello" {
		t.Errorf("Text = %q", m.Text)
	}
	if m.Format != botmsg.FormatHTML {
		t.Errorf("Format = %v, want HTML", m.Format)
	}
	if m.Keyboard != nil {
		t.Errorf("expected no keyboard")
	}
	if m.IsEdit {
		t.Errorf("expected not an edit")
	}
}

func TestRenderPromptGrid(t *testing.T) {
	plan := botplan.MessagePlan{
		Text: botplan.RichText("Coming?"),
		Prompt: &botplan.ActionPrompt{
			LayoutRows: 2,
			Choices: []botplan.Choice{
				{Label: "Yes", Token: "rsvp_y_1"},
				{Label: "Maybe", Token: "rsvp_m_1"},
				{Label: "No", Token: "rsvp_n_1"},
				{Label: "Later", Token: "rsvp_l_1"},
			},
		},
	}
	msgs, err := NewRenderer().Render(plan, tgTarget)
	if err != nil {
		t.Fatal(err)
	}
	kb := requireInlineKeyboard(t, msgs[0].Keyboard)
	// 4 choices into 2 rows → 2 buttons per row.
	if len(kb.Buttons) != 2 {
		t.Fatalf("want 2 rows, got %d", len(kb.Buttons))
	}
	for i, row := range kb.Buttons {
		if len(row) != 2 {
			t.Fatalf("row %d has %d buttons, want 2", i, len(row))
		}
	}
	// First button's token must ride as callback_data.
	first, ok := kb.Buttons[0][0].(*botkb.DataButton)
	if !ok {
		t.Fatalf("want DataButton, got %T", kb.Buttons[0][0])
	}
	if first.Data != "rsvp_y_1" || first.Text != "Yes" {
		t.Errorf("first button = %+v", first)
	}
}

func TestRenderPromptDefaultLayoutOnePerRow(t *testing.T) {
	plan := botplan.MessagePlan{
		Text: botplan.RichText("x"),
		Prompt: &botplan.ActionPrompt{
			Choices: []botplan.Choice{
				{Label: "A", Token: "a"},
				{Label: "B", Token: "b"},
				{Label: "C", Token: "c"},
			},
		},
	}
	msgs, _ := NewRenderer().Render(plan, tgTarget)
	kb := requireInlineKeyboard(t, msgs[0].Keyboard)
	if len(kb.Buttons) != 3 {
		t.Fatalf("want 3 rows (one per button), got %d", len(kb.Buttons))
	}
}

func TestRenderURLActionButton(t *testing.T) {
	plan := botplan.MessagePlan{
		Text:      botplan.RichText("Open it"),
		URLAction: &botplan.URLAction{Label: "View day", URL: "https://x.io/d/1"},
	}
	msgs, err := NewRenderer().Render(plan, tgTarget)
	if err != nil {
		t.Fatal(err)
	}
	kb := requireInlineKeyboard(t, msgs[0].Keyboard)
	last := kb.Buttons[len(kb.Buttons)-1]
	urlBtn, ok := last[0].(*botkb.UrlButton)
	if !ok {
		t.Fatalf("want UrlButton, got %T", last[0])
	}
	if urlBtn.URL != "https://x.io/d/1" || urlBtn.Text != "View day" {
		t.Errorf("url button = %+v", urlBtn)
	}
}

func TestRenderPromptPlusURLActionInSameKeyboard(t *testing.T) {
	plan := botplan.MessagePlan{
		Text:      botplan.RichText("x"),
		Prompt:    &botplan.ActionPrompt{Choices: []botplan.Choice{{Label: "Yes", Token: "y"}}},
		URLAction: &botplan.URLAction{Label: "View", URL: "https://x.io"},
	}
	msgs, err := NewRenderer().Render(plan, tgTarget)
	if err != nil {
		t.Fatal(err)
	}
	// Telegram supports both in one message.
	if len(msgs) != 1 {
		t.Fatalf("want 1 message, got %d", len(msgs))
	}
	kb := requireInlineKeyboard(t, msgs[0].Keyboard)
	if len(kb.Buttons) != 2 { // 1 choice row + 1 url row
		t.Fatalf("want 2 rows, got %d", len(kb.Buttons))
	}
	if _, ok := kb.Buttons[0][0].(*botkb.DataButton); !ok {
		t.Error("row 0 should be the data button")
	}
	if _, ok := kb.Buttons[1][0].(*botkb.UrlButton); !ok {
		t.Error("row 1 should be the url button")
	}
}

func TestRenderLivePanelSetsEdit(t *testing.T) {
	plan := botplan.MessagePlan{
		Text:      botplan.RichText("updated"),
		LivePanel: &botplan.LivePanel{PanelKey: "card1"},
	}
	msgs, err := NewRenderer().Render(plan, tgTarget)
	if err != nil {
		t.Fatal(err)
	}
	if !msgs[0].IsEdit {
		t.Error("live panel should set IsEdit")
	}
}

func TestRenderProactiveIsOrdinarySend(t *testing.T) {
	// Telegram has no window: a proactive plan renders exactly like a reply.
	plan := botplan.MessagePlan{
		Text:      botplan.RichText("Sat kitesurf 15:00"),
		Proactive: &botplan.ProactiveSpec{Purpose: "intent_notice", Locale: "en"},
	}
	msgs, err := NewRenderer().Render(plan, tgTarget)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 || msgs[0].Text != "Sat kitesurf 15:00" || msgs[0].IsEdit {
		t.Errorf("proactive not rendered as ordinary send: %+v", msgs)
	}
}

func TestRenderMediaYieldsPhotoThenText(t *testing.T) {
	plan := botplan.MessagePlan{
		Text:  botplan.RichText("caption text below"),
		Media: &botplan.MediaRef{ImageURL: "https://x.io/a.jpg", Caption: "cap"},
	}
	msgs, err := NewRenderer().Render(plan, tgTarget)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 {
		t.Fatalf("want 2 messages (photo + text), got %d", len(msgs))
	}
	photo, ok := msgs[0].BotMessage.(SendPhoto)
	if !ok {
		t.Fatalf("first message should be SendPhoto, got %T", msgs[0].BotMessage)
	}
	if photo.Photo.(tgbotapi.PhotoUrl) != "https://x.io/a.jpg" {
		t.Errorf("photo url = %v", photo.Photo)
	}
	if photo.Caption != "cap" {
		t.Errorf("caption = %q", photo.Caption)
	}
}

func TestRenderMediaWithMediaIDUsesFileID(t *testing.T) {
	plan := botplan.MessagePlan{
		Text:  botplan.RichText("x"),
		Media: &botplan.MediaRef{MediaID: "AgACfileid"},
	}
	msgs, _ := NewRenderer().Render(plan, tgTarget)
	photo := msgs[0].BotMessage.(SendPhoto)
	if photo.Photo.(tgbotapi.FileID) != "AgACfileid" {
		t.Errorf("want FileID, got %v", photo.Photo)
	}
}

func TestRenderRejectsWrongPlatform(t *testing.T) {
	plan := botplan.MessagePlan{Text: botplan.RichText("x")}
	_, err := NewRenderer().Render(plan, botplan.RenderTarget{Platform: botplan.PlatformWhatsApp})
	if !errors.Is(err, botplan.ErrUnsupportedTarget) {
		t.Errorf("want ErrUnsupportedTarget, got %v", err)
	}
}

func TestRenderPropagatesInvalidPlan(t *testing.T) {
	plan := botplan.MessagePlan{} // no text, no media
	_, err := NewRenderer().Render(plan, tgTarget)
	if !errors.Is(err, botplan.ErrInvalidPlan) {
		t.Errorf("want ErrInvalidPlan, got %v", err)
	}
}

func TestTelegramDescriptorMatchesCapabilityFacts(t *testing.T) {
	d := NewRenderer().Descriptor()
	if !d.SupportsEdit || !d.SupportsDelete || !d.SupportsCallbackAck {
		t.Error("telegram supports edit/delete/ack")
	}
	if d.WindowGated {
		t.Error("telegram has no send window")
	}
	if !d.SupportsButtonGrid || !d.SupportsInlineURLButton || !d.SupportsMedia {
		t.Error("telegram supports grid/url-button/media")
	}
	if d.TextMarkup != botplan.MarkupHTML || !d.SupportsAnchorTextLinks {
		t.Error("telegram uses HTML with anchor links")
	}
}

func requireInlineKeyboard(t *testing.T, kb botkb.Keyboard) *botkb.MessageKeyboard {
	t.Helper()
	if kb == nil {
		t.Fatal("expected a keyboard, got nil")
	}
	mk, ok := kb.(*botkb.MessageKeyboard)
	if !ok {
		t.Fatalf("want *botkb.MessageKeyboard, got %T", kb)
	}
	if mk.KeyboardType() != botkb.KeyboardTypeInline {
		t.Fatalf("want inline keyboard, got %v", mk.KeyboardType())
	}
	return mk
}
