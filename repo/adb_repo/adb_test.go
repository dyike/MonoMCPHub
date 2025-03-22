package adb_repo

import "testing"

var adbRepo AdbRepo

func init() {
	adbRepo = NewAdbRepo("192.168.5.72:33479", "./")
}
func TestGetPackages(t *testing.T) {
	packages, err := adbRepo.GetPackages()
	if err != nil {
		t.Fatalf("Failed to get packages: %v", err)
	}
	t.Logf("Packages: %v", packages)
}

func TestGetScreenshot(t *testing.T) {
	err := adbRepo.TakeScreenshot()
	if err != nil {
		t.Fatalf("Failed to take screenshot: %v", err)
	}
}

func TestGetUILayout(t *testing.T) {
	uiLayout, err := adbRepo.GetUILayout()
	if err != nil {
		t.Fatalf("Failed to get ui layout: %v", err)
	}
	t.Logf("UILayout: %v", uiLayout)
}
