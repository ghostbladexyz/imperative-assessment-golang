package receipts

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

func TestIssueAndValidate(t *testing.T) {
	t.Parallel()
	manager, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := manager.Issue("foundation/3", "source-hash")
	if err != nil {
		t.Fatal(err)
	}
	receipt, valid := manager.Validate(encoded)
	if !valid || receipt.ExerciseKey != "foundation/3" || receipt.SourceHash != "source-hash" {
		t.Fatalf("unexpected receipt: %#v, valid=%v", receipt, valid)
	}
	if _, valid := manager.Validate(encoded + "tampered"); valid {
		t.Fatal("tampered receipt passed validation")
	}
}

func TestValidateLegacyPositionReceipt(t *testing.T) {
	t.Parallel()
	manager, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	legacy := Receipt{LevelID: 171, SourceHash: "legacy-source", IssuedAt: time.Now().Unix()}
	legacy.Signature = manager.signLegacyPosition(legacy.LevelID, legacy.SourceHash, legacy.IssuedAt)
	data, err := json.Marshal(legacy)
	if err != nil {
		t.Fatal(err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(data)
	receipt, valid := manager.Validate(encoded)
	if !valid || receipt.LevelID != 171 || receipt.ExerciseKey != assessment.ExerciseKey("") {
		t.Fatalf("unexpected legacy receipt: %#v, valid=%v", receipt, valid)
	}
}
