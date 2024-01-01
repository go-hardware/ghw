//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package memory_test

import (
	"encoding/json"
	"testing"

	ghwcontext "github.com/go-hardware/ghw/pkg/context"
	"github.com/go-hardware/ghw/pkg/memory"
)

// we have this test in memory_linux_test.go (and not in memory_test.go) because `mem.load.Info` is implemented
// only on linux; so having it in the platform-independent tests would lead to false negatives.
func TestMemoryMarshalUnmarshal(t *testing.T) {
	ctx := ghwcontext.New(ghwcontext.WithDisableWarnings())
	data, err := memory.New(ctx)
	if err != nil {
		t.Fatalf("Expected no error creating memory.Info, but got %v", err)
	}

	jdata, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Expected no error marshaling memory.Info, but got %v", err)
	}

	var topo *memory.Info

	err = json.Unmarshal(jdata, &topo)
	if err != nil {
		t.Fatalf("Expected no error unmarshaling memory.Info, but got %v", err)
	}
}
