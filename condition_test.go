package main

import (
	"testing"
)

func TestValid_OriginSelf_OtherCreatureReturnsFalse(t *testing.T) {
	c := &Condition{"origin", []string{"self"}}

	creature1 := NewTestEntity("Dodgy Fella", "p1")
	creature2 := NewTestEntity("Dodgy Fella", "p1")

	if c.Valid(creature1, creature2) {
		t.Errorf("Expected other creature to be invalid, but was valid for origin self")
	}
}

func TestValid_OriginSelf_SameCreatureReturnsTrue(t *testing.T) {
	c := &Condition{"origin", []string{"self"}}

	creature1 := NewTestEntity("Dodgy Fella", "p1")

	if c.Valid(creature1, creature1) == false {
		t.Errorf("Expected same creature to be valid, but was invalid for origin self")
	}
}

func TestValid_OriginOther_SameCreatureShouldReturnFalse(t *testing.T) {
	c := &Condition{"origin", []string{"other"}}

	creature1 := NewTestEntity("Dodgy Fella", "p1")

	if c.Valid(creature1, creature1) {
		t.Errorf("Expected same creature to be invalid, but was valid for origin other")
	}
}

func TestValid_OriginOther_OtherCreatureShouldReturnTrue(t *testing.T) {
	c := &Condition{"origin", []string{"other"}}

	creature1 := NewTestEntity("Dodgy Fella", "p1")
	creature2 := NewTestEntity("Dodgy Fella", "p1")

	if c.Valid(creature1, creature2) == false {
		t.Errorf("Expected same creature to be valid, but was invalid for origin other")
	}
}
