package banner

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestPadRight(t *testing.T) {
	if s := padRight("hello", 10); s != "hello     " {
		t.Fatalf("expected 'hello     ', got %q", s)
	}
	if s := padRight("hello", 3); s != "hello" {
		t.Fatalf("expected 'hello', got %q", s)
	}
	if s := padRight("", 5); s != "     " {
		t.Fatalf("expected 5 spaces, got %q", s)
	}
}

func TestPadRightUnicode(t *testing.T) {
	s := padRight("héllo", 10)
	if utf8.RuneCountInString(s) != 10 {
		t.Fatalf("expected 10 runes, got %d for %q", utf8.RuneCountInString(s), s)
	}
	if !strings.HasPrefix(s, "héllo") {
		t.Fatalf("expected padded unicode, got %q", s)
	}
}

func TestLine(t *testing.T) {
	result := line("hello", "world")
	if !strings.HasPrefix(result, "│") {
		t.Fatal("expected line to start with │")
	}
	if !strings.HasSuffix(result, "│") {
		t.Fatal("expected line to end with │")
	}
}

func TestString(t *testing.T) {
	s := String()
	if !strings.HasPrefix(s, "╭") {
		t.Fatal("expected banner to start with ╭")
	}
	if !strings.HasSuffix(s, "╯") {
		t.Fatal("expected banner to end with ╯")
	}
	if !strings.Contains(s, "Author : RK Riad Khan") {
		t.Fatal("expected author line")
	}
	if !strings.Contains(s, "GitHub : rkriad585/neoarc") {
		t.Fatal("expected GitHub line")
	}
	if !strings.Contains(s, "Version:") {
		t.Fatal("expected version line")
	}
	if !strings.Contains(s, "Commit :") {
		t.Fatal("expected commit line")
	}
}

func TestStringWidth(t *testing.T) {
	s := String()
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		if utf8.RuneCountInString(l) != width {
			t.Fatalf("expected line width %d runes, got %d for line %q", width, utf8.RuneCountInString(l), l)
		}
	}
}

func TestPrint(t *testing.T) {
	Print()
}
