package models

import (
	"fmt"
	"os"
	"testing"
)

func respWith(profileId string) *RegistrationInfoResponse {
	return &RegistrationInfoResponse{Data: RegistrationInfo{ProfileId: profileId}}
}

func TestGetNativeAppProfile_Valid(t *testing.T) {
	p, err := GetNativeAppProfile(respWith("12345678-90ab-cdef-1234-567890abcdef"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantIpc := fmt.Sprintf("mozeidon_native_app_%d_12345678", os.Getpid())
	if p.IpcName != wantIpc {
		t.Fatalf("IpcName = %q, want %q", p.IpcName, wantIpc)
	}
	wantFile := fmt.Sprintf("%d_12345678.json", os.Getpid())
	if p.FileName != wantFile {
		t.Fatalf("FileName = %q, want %q", p.FileName, wantFile)
	}
}

func TestGetNativeAppProfile_ShortProfileId(t *testing.T) {
	if _, err := GetNativeAppProfile(respWith("abc")); err == nil {
		t.Fatal("expected error for short profileId")
	}
}

func TestGetNativeAppProfile_EmptyProfileId(t *testing.T) {
	if _, err := GetNativeAppProfile(respWith("")); err == nil {
		t.Fatal("expected error for empty profileId")
	}
}
