package statements

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/invopop/jsonschema"
	"github.com/kushturner/finances/internal/transactions"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type LLMClient struct {
	client *openai.Client
}

type transactionSchema struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Out         int64     `json:"out"`
	In          int64     `json:"in"`
	Currency    string    `json:"currency"`
}

type statementSchema struct {
	Transactions []transactionSchema `json:"transactions"`
}

func GenerateSchema[T any]() map[string]any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
		Anonymous:                 true,
	}

	var v T
	schema := reflector.Reflect(v)

	b, err := json.Marshal(schema)
	if err != nil {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		return map[string]any{}
	}
	return out
}

var StatementsResponseSchema = GenerateSchema[statementSchema]()

var schemaParams = responses.ResponseFormatTextJSONSchemaConfigParam{
	Name:        "statement",
	Description: openai.String("bank statement transactions"),
	Schema:      StatementsResponseSchema,
	Strict:      openai.Bool(true),
}

func (c *LLMClient) ReadStatement(ctx context.Context, base64Statement string) ([]transactions.Transaction, error) {
	prompt := `Extract all transactions from the following bank statement. For each transaction, provide the date using ISO 8601 format (use the primary transaction date shown in the Date column, NOT the Effective Date), description (include the whole description), amount out in its smallest unit (0 if empty), amount in in its smallest unit (0 if empty) and currency.`

	response, err := c.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: openai.ChatModelGPT4o,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemUnionParam{
					OfInputMessage: &responses.ResponseInputItemMessageParam{
						Role: "user",
						Content: responses.ResponseInputMessageContentListParam{
							responses.ResponseInputContentUnionParam{
								OfInputFile: &responses.ResponseInputFileParam{
									FileData: openai.String("data:application/pdf;base64," + base64Statement),
									Filename: openai.String("statement.pdf"),
									Type:     "input_file",
								},
							},
							responses.ResponseInputContentUnionParam{
								OfInputText: &responses.ResponseInputTextParam{
									Text: prompt,
									Type: "input_text",
								},
							},
						},
					},
				},
			},
		},
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &schemaParams,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get completion from LLM: %w", err)
	}

	var statements statementSchema
	err = json.Unmarshal([]byte(response.OutputText()), &statements)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	transactions := make([]transactions.Transaction, len(statements.Transactions))
	for i, tx := range statements.Transactions {
		transactions[i] = toTransaction(tx)
	}

	return transactions, nil
}

func NewLLMClient(apiKey string) *LLMClient {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &LLMClient{client: &client}
}

func toTransaction(tx transactionSchema) transactions.Transaction {
	return transactions.Transaction{
		Date:        tx.Date,
		Description: tx.Description,
		Out:         money.New(tx.Out, tx.Currency),
		In:          money.New(tx.In, tx.Currency),
	}
}
