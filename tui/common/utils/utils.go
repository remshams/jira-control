package tui_common_utils

import (
	"fmt"

	"github.com/remshams/common/tui/styles"
)

func RenderKeyValue(key string, value string) string {
	return fmt.Sprintf("%s%s %s", styles.TextAccentColor.Render(key), styles.TextAccentColor.Render(":"), value)
}
