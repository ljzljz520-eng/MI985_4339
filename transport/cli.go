package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"example.com/nursery-cms/domain"
	"example.com/nursery-cms/service"
)

type CLI struct{ Service *service.Service }

func NewCLI(s *service.Service) CLI { return CLI{Service: s} }

func (c CLI) Run(ctx context.Context, args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required")
	}
	switch args[0] {
	case "list":
		batch := ""
		if len(args) > 1 {
			batch = args[1]
		}
		items, err := c.Service.Search(ctx, batch, "")
		if err != nil {
			return err
		}
		return json.NewEncoder(out).Encode(items)
	case "import":
		if len(args) < 2 {
			return fmt.Errorf("batch is required")
		}
		item := domain.ImportedRecord{ExternalID: "cli-1", BatchID: args[1], Title: "课堂活动", Content: strings.Join(args[2:], " "), Owner: "cli"}
		result, err := c.Service.ImportBatch(ctx, args[1], []domain.ImportedRecord{item})
		if err != nil {
			return err
		}
		return json.NewEncoder(out).Encode(result)
	case "help":
		_, err := io.WriteString(out, "list [batch]\nimport <batch> <content>\n")
		return err
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func ParseArgs(line string) []string {
	parts := strings.Fields(line)
	return append([]string(nil), parts...)
}

func Usage() string { return "nursery-cms list [batch] | import <batch> <content>" }
