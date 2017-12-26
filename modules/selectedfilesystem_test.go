package modules

import (
	"testing"
)

func TestcleanPathName(t *testing.T) {
	// Clean name
	if "main.go" != cleanPathName("main.go") {
		t.Error("Error while cleaning clean name")
	}
	// Leading slash
	if "main.go" != cleanPathName("/main.go") {
		t.Error("Error while cleaning leading slash")
	}
	// Multiple leading slashes
	if "main.go" != cleanPathName("//////////main.go") {
		t.Error("Error while cleaning multiple leading slashes")
	}
	// Slash in name
	if "main//.go" != cleanPathName("main//.go") {
		t.Error("Error while cleaning slash in name")
	}
}

func TestNewSelectedFileSystem(t *testing.T) {
	for _, value := range []int{10, 1000, 389} {
		sfs := NewSelectedFileSystem(value)
		if sfs.files == nil {
			t.Log("Newly created SelectedFileSystem has nil file array - size", value)
			t.FailNow()
		}
		if len(sfs.files) != 0 {
			t.Log("Newly created SelectedFileSystem has non-zero file length - size", value)
			t.FailNow()
		}
		if cap(sfs.files) != value {
			t.Log("Newly created SelectedFileSystem has wrong file capacity - size", value, "capacity", cap(sfs.files))
			t.FailNow()
		}
	}

}

func TestAddFile(t *testing.T) {
	sfs := NewSelectedFileSystem(10)

	// Add 1 file
	err := sfs.AddFile("selectedfilesystem_test.go")
	if err != nil {
		t.Error(err)
	}
	if !stringSliceComparison(sfs.files, []string{cleanPathName("selectedfilesystem_test.go")}) {
		t.Error("files are not matching after 1 add")
	}

	// Add more files
	err = sfs.AddFile("selectedfilesystem_test.go")
	if err != nil {
		t.Error(err)
	}
	err = sfs.AddFile("//selectedfilesystem_test.go")
	if err != nil {
		t.Error(err)
	}
	err = sfs.AddFile("./selectedfilesystem_test.go")
	if err != nil {
		t.Error(err)
	}
	err = sfs.AddFile("selectedfilesystem_test.go")
	if err != nil {
		t.Error(err)
	}
	if !stringSliceComparison(sfs.files, []string{cleanPathName("selectedfilesystem_test.go"), cleanPathName("selectedfilesystem_test.go"), cleanPathName("//selectedfilesystem_test.go"), cleanPathName("./selectedfilesystem_test.go"), cleanPathName("selectedfilesystem_test.go")}) {
		t.Error("files are not matching after 5 adds")
	}

	// Test non-existent file
	err = sfs.AddFile("surely_not_existing")
	if err == nil {
		t.Error("Test fail: No error for non existing file")
	}

}

func TestOpen(t *testing.T) {
	sfs := NewSelectedFileSystem(10)
	sfs.AddFile("selectedfilesystem_test.go")

	// Test existing file
	file, err := sfs.Open("selectedfilesystem_test.go")
	if err != nil {
		t.Error("Error is not nil for existing file")
	}
	if file == nil {
		t.Error("File is nil for existing file")
	}

	// Test non-existing file
	file, err = sfs.Open("surely_not_existing")
	if err == nil {
		t.Error("Error is nil for non-existing file")
	}
	if file != nil {
		t.Error("File is not nil for non-existing file")
	}
}

// Returns true if both slices have the same length and the same content.
func stringSliceComparison(A, B []string) bool {
	if len(A) != len(B) {
		return false
	}
	for i := range A {
		if A[i] != B[i] {
			return false
		}
	}
	return true
}
