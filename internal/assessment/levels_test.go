package assessment

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"testing"
)

func TestCheckpointCatalogueIsCompleteAndRanked(t *testing.T) {
	want := []string{"validate-stack", "safe-sum", "tetris", "method-routing", "status-matrix", "wordcount", "ls-format", "ascii-render", "reloaded-rev", "reloaded-format", "lemin-path", "lemin-why", "bounded-fanout", "consume-join", "push-swap", "broadcast", "reactions"}
	if err := Validate(); err != nil {
		t.Fatal(err)
	}
	levels := Levels()
	if len(levels) != len(want) {
		t.Fatalf("got %d exercises, want %d", len(levels), len(want))
	}
	for index, slug := range want {
		level := levels[index]
		if level.Key != ExerciseKey("checkpoint/"+slug) || level.Title != slug || level.ID != index+1 || level.TrackPosition != index+1 || level.Track != TrackCore {
			t.Errorf("rank %d = %#v, want %q", index+1, level, slug)
		}
		if len(level.Tests) == 0 {
			t.Errorf("exercise %q is missing its grader", slug)
		}
	}
}

func TestCheckpointAssetsMatchPinnedUpstreamBlobs(t *testing.T) {
	if checkpointUpstreamCommit != "9ef45ee4a164b6d9b8ff3c5ad9b97753c3d63297" {
		t.Fatal("unexpected upstream commit")
	}
	want := map[string]string{
		"ascii-render/README.md": "824364c36d43ec9577168d79c70fba74c52748a9", "ascii-render/starter.go.txt": "1b332c788eccc234b3583dd71b681c20c4276d0a", "ascii-render/resources/banner.txt": "77f3ea44e1ef356a71ad2891e68c365923945245",
		"bounded-fanout/README.md": "0c67963c4b8e2a138291ca03c6682ed7cdec9fdf", "bounded-fanout/starter.go.txt": "593ddd303c8efc27837442598ef412f2a5b9d972",
		"broadcast/README.md": "04c01b377d703646df0b1141ad8bcdf226ea9a26", "broadcast/starter.go.txt": "56274f691cf0fa3a2b4346159934b94e3d8d53d0",
		"consume-join/README.md": "7f56b6e6aee2ea0faba1a55a02c9326c754c23d4", "consume-join/starter.go.txt": "497d616ab3b1e0be7c3c5b30e4aa7bbf2baf29f6",
		"lemin-path/README.md": "4d35c751b02209e2e84d566a6dbf0aa2d2c4c1f8", "lemin-path/starter.go.txt": "4903c7fd8b75627c4aa3bcdd2ad12501f735dabb",
		"lemin-why/README.md": "f31af500822376bc68c281574c2ab1eabfd17c04", "lemin-why/starter.go.txt": "271148be649177c6a1feaf7db836608e9c93c466",
		"ls-format/README.md": "8d80261d0d7fe231eb9ac0ca673e522e6d7d91f6", "ls-format/starter.go.txt": "6ecca243e989b4996d5e127c7c164565daa1cb9e",
		"method-routing/README.md": "420570b46f298f83660c647143e5edb11daf36e1", "method-routing/starter.go.txt": "8e1b28eac3c72acd6d56707236b524c280fb23c0",
		"push-swap/README.md": "934708ed63423ed00ec63b8955b143d5808d8d45", "push-swap/starter.go.txt": "fc45d99340f4784b06a1ed3f84462f5b850f8922",
		"reactions/README.md": "7ce2c9ce5d55ff8d0f599e8380a85af9822a21ca", "reactions/starter.go.txt": "da4b6533868eb8923ebc27eb3c36b3875c46a394", "reactions/resources/schema.sql": "f9df1378ab73f8f23fe9ac8f9cc5e436db01cebc",
		"reloaded-format/README.md": "f173f18567160a14d6b09befe5e37d4dd0e78c5f", "reloaded-format/starter.go.txt": "1914e7f5a8bd1e85e687a36993a51d5294e8c1ad",
		"reloaded-rev/README.md": "dc8bfec54f702f6ee88dae8649f86c847991acbc", "reloaded-rev/starter.go.txt": "248b83c9508485c92cc9174cf8f34f5733d05e1b",
		"safe-sum/README.md": "15b9aa85b79fe2de3df9f7eea0c068e09bf1feca", "safe-sum/starter.go.txt": "2021fb94ddf685d6cda397e433325646fc65dec7",
		"status-matrix/README.md": "5237f6d52516d8bab78c544f0ded37a4df6b3f81", "status-matrix/starter.go.txt": "1fb7eedf55a6afc1c6ec4b358a4075c11fd6f1f4",
		"tetris/README.md": "5e473f724c37fd60efc456d8b26f151db422815c", "tetris/starter.go.txt": "1aa40983ed6e6eeae54b77faa394c7b53ee55f9d",
		"validate-stack/README.md": "17224ff09cbf37cb72bac3e44df2a2655c9a9c9e", "validate-stack/starter.go.txt": "03c1e6368b736e6194c234d7ca932a030d85d966",
		"wordcount/README.md": "12ad8256f0f7bb4b26e5a63e1fc7ef532e1ec2c9", "wordcount/starter.go.txt": "1fb9ace3142d2d7b98ff32bd1d4c6fb1029733e8",
	}
	for path, blobID := range want {
		content, err := checkpointAssets.ReadFile("checkpoint/" + path)
		if err != nil {
			t.Fatal(err)
		}
		content = []byte(strings.ReplaceAll(string(content), "\r\n", "\n"))
		header := fmt.Sprintf("blob %d\x00", len(content))
		actual := fmt.Sprintf("%x", sha1.Sum(append([]byte(header), content...)))
		if actual != blobID {
			t.Errorf("%s blob = %s, want %s", path, actual, blobID)
		}
	}
}

func TestPublicLevelsPublishOfficialChecks(t *testing.T) {
	for _, level := range PublicLevels() {
		if level.Instructions.AllowedBuiltins == nil || level.Instructions.AllowedPackages == nil {
			t.Fatalf("%q serializes a null allowlist", level.Key)
		}
		if level.Resources == nil {
			t.Fatalf("%q serializes a null resource list", level.Key)
		}
		for _, current := range level.Tests {
			if strings.TrimSpace(current.Input) == "" || strings.TrimSpace(current.Expected) == "" || current.Expected == "pass" {
				t.Errorf("%q test %q has placeholder input/need", level.Key, current.ID)
			}
		}
	}
}

func TestCatalogueReturnsDefensiveProjections(t *testing.T) {
	first := Levels()
	first[0].Title = "changed"
	first[0].Tests[0].Name = "changed"
	second := Levels()
	if second[0].Title == "changed" || second[0].Tests[0].Name == "changed" {
		t.Fatal("caller mutation changed cached catalogue")
	}
}
