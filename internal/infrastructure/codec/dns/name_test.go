package dns

import "testing"

func TestEncodeName(t *testing.T) {
	buf := encodeName("example.com")

	expected := []byte{
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		3, 'c', 'o', 'm',
		0,
	}

	if len(buf) != len(expected) {
		t.Fatalf("length mismatch: %d vs %d", len(buf), len(expected))
	}

	for i := range expected {
		if buf[i] != expected[i] {
			t.Fatalf("byte mismatch at %d", i)
		}
	}
}

func TestParseName_Simple(t *testing.T) {
	buf := []byte{
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		3, 'c', 'o', 'm',
		0,
	}

	name, off := parseName(buf, 0)

	if name != "example.com" {
		t.Fatalf("unexpected name: %s", name)
	}

	if off != len(buf) {
		t.Fatalf("unexpected offset: %d", off)
	}
}

func TestParseName_Compressed(t *testing.T) {
	buf := []byte{
		7, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		3, 'c', 'o', 'm',
		0,
		3, 'w', 'w', 'w',
		0xC0, 0x00,
	}

	name, off := parseName(buf, 13)

	if name != "www.example.com" {
		t.Fatalf("unexpected name: %s", name)
	}

	if off != 19 {
		t.Fatalf("unexpected offset: %d", off)
	}
}
