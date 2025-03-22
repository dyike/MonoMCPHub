package adb_repo

import "testing"

func TestGetPackages(t *testing.T) {
	adbRepo := NewAdbRepo("127.0.0.1:5555")
	packages, err := adbRepo.GetPackages()
	if err != nil {
		t.Fatalf("Failed to get packages: %v", err)
	}
	t.Logf("Packages: %v", packages)
}
