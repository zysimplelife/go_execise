package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
)

func main() {

	original := []byte(`{"name": "abc", "occupation": {"title":"test"}}`)

	// Let's create a merge patch from these two documents...
	patch1 := []byte(`{"occupation": {"title" : "title"}}`)
	patch2 := []byte(`{"occupation": {"years" : "16"}}`)
	patch3 := []byte(`{"occupation": {"heir" : "Joffrey"}}`)

	combinedPatch, err := jsonpatch.MergeMergePatches(patch1, patch2)
	if err != nil {
		panic(err)
	}

	combinedPatch, err = jsonpatch.MergeMergePatches(combinedPatch, patch3)
	if err != nil {
		panic(err)
	}

	patchJSON := fmt.Sprintf(`[{"op": "replace", "path": "/occupation", "value": %s}]`, combinedPatch)

	patch, err := jsonpatch.DecodePatch([]byte(patchJSON))
	if err != nil {
		panic(err)
	}

	modified, err := patch.Apply(original)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Patch1 document: %s\n", patch1)
	fmt.Printf("Patch2 document: %s\n", patch2)
	fmt.Printf("Patch3 document: %s\n", patch3)
	fmt.Printf("MergedPatch document: %s\n", patchJSON)
	fmt.Printf("Original document: %s\n", original)
	fmt.Printf("Modified document: %s\n", modified)

}
