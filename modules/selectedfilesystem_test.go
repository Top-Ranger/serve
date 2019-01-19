// SPDX-License-Identifier: Apache-2.0
// Copyright 2017,2018,2019 Marcus Soll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package modules

import (
	"testing"
)

func TestCleanPathName(t *testing.T) {
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

	// Test non-existent file
	err = sfs.AddFile("../modules/")
	if err == nil {
		t.Error("Test fail: No error for adding directory")
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

func TestStringSliceComparison(t *testing.T) {
	if !stringSliceComparison([]string{"cake", "cookie", "biscuit"}, []string{"cake", "cookie", "biscuit"}) {
		t.Error()
	}
	if stringSliceComparison([]string{"cake", "cookie"}, []string{"cake", "cookie", "biscuit"}) {
		t.Error()
	}
	if stringSliceComparison([]string{"cake", "cookie", "biscuit"}, []string{"cake", "cookie"}) {
		t.Error()
	}

	// Test different capacity
	A := make([]string, 0, 10)
	B := make([]string, 0, 100)
	C := make([]string, 0, 3)

	A = append(A, "cake", "cookie", "biscuit")
	B = append(B, "cake", "cookie", "biscuit")
	C = append(C, "cake", "cookie")

	if !stringSliceComparison(A, B) {
		t.Error()
	}
	if stringSliceComparison(A, C) {
		t.Error()
	}
	if !stringSliceComparison(B, A) {
		t.Error()
	}
	if stringSliceComparison(B, C) {
		t.Error()
	}
	if stringSliceComparison(C, A) {
		t.Error()
	}
	if stringSliceComparison(C, B) {
		t.Error()
	}
}

func TestGetFiles(t *testing.T) {
	sfs := NewSelectedFileSystem(10)

	// Add 1 file
	err := sfs.AddFile("selectedfilesystem_test.go")
	if err != nil {
		t.Error(err)
	}
	if !stringSliceComparison(sfs.GetFiles(), []string{"selectedfilesystem_test.go"}) {
		t.Error("files are not matching after 1 add")
	}

	// Add more files
	err = sfs.AddFile("selectedfilesystem.go")
	if err != nil {
		t.Error(err)
	}
	if !stringSliceComparison(sfs.GetFiles(), []string{"selectedfilesystem_test.go", "selectedfilesystem.go"}) {
		t.Error("files are not matching after 2 adds")
	}
}
