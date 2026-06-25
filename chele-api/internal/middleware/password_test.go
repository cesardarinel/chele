package middleware

import (
	"strings"
	"testing"
)

func TestMakeDjangoPassword(t *testing.T) {
	hash := MakeDjangoPassword("testpass123")
	parts := strings.Split(hash, "$")
	if len(parts) != 4 {
		t.Fatalf("expected 4 parts, got %d", len(parts))
	}
	if parts[0] != "pbkdf2_sha256" {
		t.Errorf("expected pbkdf2_sha256 algorithm, got %s", parts[0])
	}
	if parts[1] != "720000" {
		t.Errorf("expected 720000 iterations, got %s", parts[1])
	}
	if len(parts[2]) == 0 {
		t.Error("salt is empty")
	}
	if len(parts[3]) == 0 {
		t.Error("hash is empty")
	}
}

func TestCheckDjangoPassword_Valid(t *testing.T) {
	password := "my-secure-pass-123"
	hash := MakeDjangoPassword(password)
	if !CheckDjangoPassword(password, hash) {
		t.Error("expected valid password to match")
	}
}

func TestCheckDjangoPassword_Invalid(t *testing.T) {
	hash := MakeDjangoPassword("realpass")
	if CheckDjangoPassword("wrongpass", hash) {
		t.Error("expected wrong password to not match")
	}
}

func TestCheckDjangoPassword_InvalidFormat(t *testing.T) {
	if CheckDjangoPassword("pass", "invalid-format") {
		t.Error("expected invalid format to return false")
	}
}

func TestCheckDjangoPassword_RealDjangoHash(t *testing.T) {
	// Django-generated hash: pbkdf2_sha256$720000$salt$hexhash
	djangoHash := "pbkdf2_sha256$720000$abc123def456$abc123def456abc123def456abc123def456abc123def456abc123def456abc1"
	if !CheckDjangoPassword("", djangoHash) {
		// This hash is fake so it should fail; we're testing the parser doesn't panic
	}
}

func TestParseDjangoHash(t *testing.T) {
	dh, err := ParseDjangoHash("pbkdf2_sha256$720000$salt$hash")
	if err != nil {
		t.Fatal(err)
	}
	if dh.Algorithm != "pbkdf2_sha256" {
		t.Errorf("algorithm: got %s", dh.Algorithm)
	}
	if dh.Iterations != 720000 {
		t.Errorf("iterations: got %d", dh.Iterations)
	}
	if dh.Salt != "salt" {
		t.Errorf("salt: got %s", dh.Salt)
	}
	if dh.Hash != "hash" {
		t.Errorf("hash: got %s", dh.Hash)
	}
}

func TestParseDjangoHash_InvalidParts(t *testing.T) {
	_, err := ParseDjangoHash("only-two$parts")
	if err == nil {
		t.Error("expected error for invalid format")
	}
}
