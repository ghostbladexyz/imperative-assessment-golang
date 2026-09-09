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
		"ascii-render/README.md": "907f0937714e053ef1fddd21cb2c348ca8b7d13c", "ascii-render/starter.go.txt": "1b332c788eccc234b3583dd71b681c20c4276d0a", "ascii-render/resources/banner.txt": "77f3ea44e1ef356a71ad2891e68c365923945245",
		"bounded-fanout/README.md": "57785d7b24820b6911f1a94ee6f7d06cc0fe274e", "bounded-fanout/starter.go.txt": "593ddd303c8efc27837442598ef412f2a5b9d972",
		"broadcast/README.md": "4998940de470adf833fbca829f77585c5aa2fc50", "broadcast/starter.go.txt": "56274f691cf0fa3a2b4346159934b94e3d8d53d0",
		"consume-join/README.md": "937056c3443198323864800ae6ed96fcf36e25c4", "consume-join/starter.go.txt": "497d616ab3b1e0be7c3c5b30e4aa7bbf2baf29f6",
		"lemin-path/README.md": "e9fd3bf93679fb2682d20ce39c6d1a9f556e6ae6", "lemin-path/starter.go.txt": "4903c7fd8b75627c4aa3bcdd2ad12501f735dabb",
		"lemin-why/README.md": "dd3174e669984a43aa66b9c006a90f0e632a67af", "lemin-why/starter.go.txt": "271148be649177c6a1feaf7db836608e9c93c466",
		"ls-format/README.md": "25cdcc75ab3b8de518f6d747e8b72ada71baecc3", "ls-format/starter.go.txt": "6ecca243e989b4996d5e127c7c164565daa1cb9e",
		"method-routing/README.md": "8746d517d1b3fcf7ec0842f22208ebd83fdceb77", "method-routing/starter.go.txt": "8e1b28eac3c72acd6d56707236b524c280fb23c0",
		"push-swap/README.md": "1f229568b790ee5ed6a8a694c60ed9454c4f78cf", "push-swap/starter.go.txt": "fc45d99340f4784b06a1ed3f84462f5b850f8922",
		"reactions/README.md": "0dce41ddfec9578103d9c3656e995b0583f86df0", "reactions/starter.go.txt": "da4b6533868eb8923ebc27eb3c36b3875c46a394", "reactions/resources/schema.sql": "f9df1378ab73f8f23fe9ac8f9cc5e436db01cebc",
		"reloaded-format/README.md": "4935cca1c266bc77c73681ffd2b4e6d713beda62", "reloaded-format/starter.go.txt": "1914e7f5a8bd1e85e687a36993a51d5294e8c1ad",
		"reloaded-rev/README.md": "a2eb0266c6237025ea4aeb3cdbec3bbd8e77fc5b", "reloaded-rev/starter.go.txt": "248b83c9508485c92cc9174cf8f34f5733d05e1b",
		"safe-sum/README.md": "6b2100ea1bc5733331630ee5233b825d572e0ab9", "safe-sum/starter.go.txt": "2021fb94ddf685d6cda397e433325646fc65dec7",
		"status-matrix/README.md": "f5ef976b136c258c430dec1a7ea7287c32ec863a", "status-matrix/starter.go.txt": "1fb7eedf55a6afc1c6ec4b358a4075c11fd6f1f4",
		"tetris/README.md": "515606a1e5f6f8689630c4e3f27ab3cd26b8d8e6", "tetris/starter.go.txt": "1aa40983ed6e6eeae54b77faa394c7b53ee55f9d",
		"validate-stack/README.md": "3c21e0454493f1ba6334571c51cfb69e65baf262", "validate-stack/starter.go.txt": "03c1e6368b736e6194c234d7ca932a030d85d966",
		"wordcount/README.md": "b705d070632b4a36061cb3fcdcf2dc6a534975ae", "wordcount/starter.go.txt": "1fb9ace3142d2d7b98ff32bd1d4c6fb1029733e8",
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
