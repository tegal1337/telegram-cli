package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zelenin/go-tdlib/client"
)

// EntitiesToMarkdown converts Telegram formatted text entities to markdown.
func EntitiesToMarkdown(text *client.FormattedText) string {
	if text == nil {
		return ""
	}

	if len(text.Entities) == 0 {
		return text.Text
	}

	runes := []rune(text.Text)

	// Sort entities by offset.
	entities := make([]*client.TextEntity, len(text.Entities))
	copy(entities, text.Entities)
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Offset < entities[j].Offset
	})

	var b strings.Builder
	lastEnd := int32(0)

	for _, e := range entities {
		// Append text before this entity.
		if e.Offset > lastEnd {
			b.WriteString(string(runes[lastEnd:e.Offset]))
		}

		entityText := string(runes[e.Offset : e.Offset+e.Length])

		switch t := e.Type.(type) {
		case *client.TextEntityTypeBold:
			b.WriteString("**")
			b.WriteString(entityText)
			b.WriteString("**")

		case *client.TextEntityTypeItalic:
			b.WriteString("*")
			b.WriteString(entityText)
			b.WriteString("*")

		case *client.TextEntityTypeUnderline:
			b.WriteString("__")
			b.WriteString(entityText)
			b.WriteString("__")

		case *client.TextEntityTypeStrikethrough:
			b.WriteString("~~")
			b.WriteString(entityText)
			b.WriteString("~~")

		case *client.TextEntityTypeCode:
			b.WriteString("`")
			b.WriteString(entityText)
			b.WriteString("`")

		case *client.TextEntityTypePre:
			b.WriteString("```\n")
			b.WriteString(entityText)
			b.WriteString("\n```")

		case *client.TextEntityTypePreCode:
			b.WriteString(fmt.Sprintf("```%s\n", t.Language))
			b.WriteString(entityText)
			b.WriteString("\n```")

		case *client.TextEntityTypeTextUrl:
			b.WriteString(fmt.Sprintf("[%s](%s)", entityText, t.Url))

		case *client.TextEntityTypeUrl:
			b.WriteString(entityText)

		case *client.TextEntityTypeMention:
			b.WriteString(entityText)

		case *client.TextEntityTypeMentionName:
			b.WriteString(fmt.Sprintf("@[%s](user:%d)", entityText, t.UserId))

		case *client.TextEntityTypeHashtag:
			b.WriteString(entityText)

		case *client.TextEntityTypeBotCommand:
			b.WriteString("`")
			b.WriteString(entityText)
			b.WriteString("`")

		case *client.TextEntityTypeEmailAddress:
			b.WriteString(entityText)

		case *client.TextEntityTypeSpoiler:
			b.WriteString("||")
			b.WriteString(entityText)
			b.WriteString("||")

		case *client.TextEntityTypeBlockQuote:
			lines := strings.Split(entityText, "\n")
			for i, line := range lines {
				b.WriteString("> ")
				b.WriteString(line)
				if i < len(lines)-1 {
					b.WriteString("\n")
				}
			}

		default:
			b.WriteString(entityText)
		}

		lastEnd = e.Offset + e.Length
	}

	// Append remaining text.
	if lastEnd < int32(len(runes)) {
		b.WriteString(string(runes[lastEnd:]))
	}

	return b.String()
}

// EntitiesToANSI converts Telegram formatted text directly to ANSI escape codes.
func EntitiesToANSI(text *client.FormattedText) string {
	if text == nil {
		return ""
	}

	if len(text.Entities) == 0 {
		return text.Text
	}

	runes := []rune(text.Text)
	entities := make([]*client.TextEntity, len(text.Entities))
	copy(entities, text.Entities)
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Offset < entities[j].Offset
	})

	var b strings.Builder
	lastEnd := int32(0)

	for _, e := range entities {
		if e.Offset > lastEnd {
			b.WriteString(string(runes[lastEnd:e.Offset]))
		}

		entityText := string(runes[e.Offset : e.Offset+e.Length])

		switch e.Type.(type) {
		case *client.TextEntityTypeBold:
			b.WriteString("\033[1m")
			b.WriteString(entityText)
			b.WriteString("\033[22m")
		case *client.TextEntityTypeItalic:
			b.WriteString("\033[3m")
			b.WriteString(entityText)
			b.WriteString("\033[23m")
		case *client.TextEntityTypeUnderline:
			b.WriteString("\033[4m")
			b.WriteString(entityText)
			b.WriteString("\033[24m")
		case *client.TextEntityTypeStrikethrough:
			b.WriteString("\033[9m")
			b.WriteString(entityText)
			b.WriteString("\033[29m")
		case *client.TextEntityTypeCode, *client.TextEntityTypePre:
			b.WriteString("\033[7m") // inverse
			b.WriteString(entityText)
			b.WriteString("\033[27m")
		default:
			b.WriteString(entityText)
		}

		lastEnd = e.Offset + e.Length
	}

	if lastEnd < int32(len(runes)) {
		b.WriteString(string(runes[lastEnd:]))
	}

	return b.String()
}
