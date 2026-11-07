package build

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

func readFrom(files map[string]string) func(string) ([]byte, error) {
	return func(name string) ([]byte, error) {
		src, ok := files[name]
		if !ok {
			return nil, fs.ErrNotExist
		}
		return []byte(src), nil
	}
}

func TestProject(t *testing.T) {
	files := map[string]string{
		"app/main.fxma":  "\"./lib/a\" aste a werer\n\"fxma:meme\" aste meme werer\n",
		"app/lib/a.fxma": "\"../b\" no f werer\n\"./c\" aste c werer\n",
		"app/b.fxma":     "oqun weho f()\nfinend\n",
		"app/lib/c.fxma": "\"../b\" aste b werer\n",
	}
	outputs, _, err := Project("app/main.fxma", readFrom(files))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, o := range outputs {
		names = append(names, o.Name)
	}
	want := []string{"app/main.fxmac", "app/lib/a.fxmac", "app/b.fxmac", "app/lib/c.fxmac"}
	if len(names) != len(want) {
		t.Fatalf("got %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("got %v, want %v", names, want)
		}
	}
	if outputs[2].Source != "Oqun We-Ho f()\nFineNd\n" {
		t.Errorf("清書 %q", outputs[2].Source)
	}
}

func TestProjectMissingImport(t *testing.T) {
	files := map[string]string{"main.fxma": "\"./nothing\" aste n werer\n"}
	_, _, err := Project("main.fxma", readFrom(files))
	var de *diag.Error
	if !errors.As(err, &de) || de.Kind != diag.NiiExi {
		t.Fatalf("エラー %v, want Nii eXi'", err)
	}
}
