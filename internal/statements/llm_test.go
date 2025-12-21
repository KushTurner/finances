package statements

import (
	"context"
	"embed"
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed *.pdf
var pdfData embed.FS

func TestReadStatement(t *testing.T) {
	t.Skip("skipping llm tests as it costs money")
	t.Run("can query llm with statement", func(t *testing.T) {
		llm := NewLLMClient(os.Getenv("OPENAPI_KEY"))

		fileName := ""
		pdf, _ := pdfData.ReadFile(fileName)
		base64PDF := base64.StdEncoding.EncodeToString(pdf)

		transaction, err := llm.ReadStatement(context.Background(), base64PDF)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		assert.Len(t, transaction, 0)
	})
}
