package receipts

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

const receiptVersion = 2

type Receipt struct {
	Version     int                    `json:"version,omitempty"`
	ExerciseKey assessment.ExerciseKey `json:"exerciseKey,omitempty"`
	LevelID     int                    `json:"levelId,omitempty"`
	SourceHash  string                 `json:"sourceHash"`
	IssuedAt    int64                  `json:"issuedAt"`
	Signature   string                 `json:"signature"`
}

type Manager struct {
	key []byte
}

func New(dataDir string) (*Manager, error) {
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return nil, err
	}
	path := filepath.Join(dataDir, "receipt.key")
	key, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		key = make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, key, 0o600); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	if len(key) < 32 {
		return nil, fmt.Errorf("receipt key is invalid")
	}
	return &Manager{key: key}, nil
}

func (manager *Manager) Issue(exerciseKey assessment.ExerciseKey, sourceHash string) (string, error) {
	receipt := Receipt{
		Version: receiptVersion, ExerciseKey: exerciseKey,
		SourceHash: sourceHash, IssuedAt: time.Now().Unix(),
	}
	receipt.Signature = manager.signExercise(receipt.ExerciseKey, receipt.SourceHash, receipt.IssuedAt)
	encoded, err := json.Marshal(receipt)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encoded), nil
}

func (manager *Manager) Validate(encoded string) (Receipt, bool) {
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Receipt{}, false
	}
	var receipt Receipt
	if json.Unmarshal(data, &receipt) != nil {
		return Receipt{}, false
	}
	var expected string
	switch {
	case receipt.Version == receiptVersion && receipt.ExerciseKey != "" && receipt.LevelID == 0:
		expected = manager.signExercise(receipt.ExerciseKey, receipt.SourceHash, receipt.IssuedAt)
	case receipt.Version == 0 && receipt.ExerciseKey == "" && receipt.LevelID > 0:
		expected = manager.signLegacyPosition(receipt.LevelID, receipt.SourceHash, receipt.IssuedAt)
	default:
		return Receipt{}, false
	}
	return receipt, hmac.Equal([]byte(expected), []byte(receipt.Signature))
}

func (manager *Manager) signExercise(exerciseKey assessment.ExerciseKey, sourceHash string, issuedAt int64) string {
	hash := hmac.New(sha256.New, manager.key)
	fmt.Fprintf(hash, "%d:%s:%s:%d", receiptVersion, exerciseKey, sourceHash, issuedAt)
	return hex.EncodeToString(hash.Sum(nil))
}

func (manager *Manager) signLegacyPosition(levelID int, sourceHash string, issuedAt int64) string {
	hash := hmac.New(sha256.New, manager.key)
	fmt.Fprintf(hash, "%d:%s:%d", levelID, sourceHash, issuedAt)
	return hex.EncodeToString(hash.Sum(nil))
}
