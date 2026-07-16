package telegram

import (
	"strings"

	"github.com/bots-go-framework/bots-fw/botplan"
)

// richToHTML renders a botplan.Rich as Telegram HTML (parse_mode=HTML).
//
// Telegram's HTML dialect supports <b>, <i>, <code>, <a href> and <blockquote>
// (capability-map telegram/text-formatting constraints.htmlTags), so every span
// kind maps to a real tag and nothing is lost. Link spans keep their anchor text
// over a hidden URL — the affordance WhatsApp lacks.
//
// All literal text is HTML-escaped so a user-supplied "<" never becomes markup.
func richToHTML(r botplan.Rich) string {
	var sb strings.Builder
	for i, line := range r.Lines {
		if i > 0 {
			sb.WriteByte('\n')
		}
		switch line.Kind {
		case botplan.LineListItem:
			sb.WriteString("• ")
			writeSpansHTML(&sb, line.Spans)
		case botplan.LineQuote:
			sb.WriteString("<blockquote>")
			writeSpansHTML(&sb, line.Spans)
			sb.WriteString("</blockquote>")
		default:
			writeSpansHTML(&sb, line.Spans)
		}
	}
	return sb.String()
}

func writeSpansHTML(sb *strings.Builder, spans []botplan.Span) {
	for _, sp := range spans {
		switch sp.Kind {
		case botplan.SpanBold:
			sb.WriteString("<b>")
			sb.WriteString(escapeHTML(sp.Text))
			sb.WriteString("</b>")
		case botplan.SpanItalic:
			sb.WriteString("<i>")
			sb.WriteString(escapeHTML(sp.Text))
			sb.WriteString("</i>")
		case botplan.SpanCode:
			sb.WriteString("<code>")
			sb.WriteString(escapeHTML(sp.Text))
			sb.WriteString("</code>")
		case botplan.SpanLink:
			sb.WriteString(`<a href="`)
			sb.WriteString(escapeHTMLAttr(sp.URL))
			sb.WriteString(`">`)
			anchor := sp.Text
			if anchor == "" {
				anchor = sp.URL
			}
			sb.WriteString(escapeHTML(anchor))
			sb.WriteString("</a>")
		default:
			sb.WriteString(escapeHTML(sp.Text))
		}
	}
}

// escapeHTML escapes the characters Telegram requires escaped in HTML text
// content: &, < and >. Quotes are not special in text content.
// https://core.telegram.org/bots/api#html-style
func escapeHTML(s string) string {
	return htmlTextEscaper.Replace(s)
}

// escapeHTMLAttr escapes for an attribute value (the href), where the quote must
// also be escaped so it cannot close the attribute early.
func escapeHTMLAttr(s string) string {
	return htmlAttrEscaper.Replace(s)
}

var (
	htmlTextEscaper = strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
	htmlAttrEscaper = strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
	)
)
