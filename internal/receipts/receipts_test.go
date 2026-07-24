package receipts

import (
	"testing"
)

func TestIssueAndValidate(t *testing.T) {
	t.Parallel()
	manager, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := manager.Issue(3, "source-hash")
	if err != nil {
		t.Fatal(err)
	}
	receipt, valid := manager.Validate(encoded)
	if !valid || receipt.LevelID != 3 || receipt.SourceHash != "source-hash" {
		t.Fatalf("unexpected receipt: %#v, valid=%v", receipt, valid)
	}
	if _, valid := manager.Validate(encoded + "tampered"); valid {
		t.Fatal("tampered receipt passed validation")
	}
}
