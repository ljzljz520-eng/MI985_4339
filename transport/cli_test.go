package transport

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"example.com/nursery-cms/service"
	"example.com/nursery-cms/store"
)

func TestCLICommands(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/cli.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := service.New(db)
	cli := NewCLI(s)
	var out bytes.Buffer
	if err := cli.Run(context.Background(), []string{"help"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "import") {
		t.Fatal(out.String())
	}
	out.Reset()
	if err := cli.Run(context.Background(), []string{"import", "batch-cli", "活动内容"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "batch-cli") {
		t.Fatal(out.String())
	}
}
