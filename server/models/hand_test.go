package models

import (
	"reflect"
	"testing"
)

func TestPickableCardWhenTableIsEmpty(t *testing.T) {
	var table Table
	var hand Hand
	var card Card

	hand.cards = []Card{
		card.fromId(Parrot2),
		card.fromId(Roger5),
		card.fromId(SkullKing),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Parrot2,
		Roger5,
		SkullKing,
		Map3,
		Map2,
	}

	options := hand.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}

func TestPickableCardWhenFirstCardOnTableIsSuit(t *testing.T) {
	var table Table
	var hand Hand
	var card Card

	table.cards = []Card{
		card.fromId(Parrot3),
	}

	hand.cards = []Card{
		card.fromId(Parrot2),
		card.fromId(Roger5),
		card.fromId(SkullKing),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Parrot2,
		SkullKing,
	}

	options := hand.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}

func TestUserCanPickAnyCardIfNoCardMatchesTheSuit(t *testing.T) {
	var table Table
	var hand Hand
	var card Card

	table.cards = []Card{
		card.fromId(Parrot3),
	}

	hand.cards = []Card{
		card.fromId(Chest1),
		card.fromId(Roger5),
		card.fromId(Pirate1),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Chest1,
		Roger5,
		Pirate1,
		Map3,
		Map2,
	}

	options := hand.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}

func TestPickableCardWhenFirstCardOnTableIsEscape1(t *testing.T) {
	var table Table
	var hand Hand
	var card Card

	table.cards = []Card{
		card.fromId(Escape1),
	}

	hand.cards = []Card{
		card.fromId(Chest1),
		card.fromId(Roger5),
		card.fromId(Pirate1),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Chest1,
		Roger5,
		Pirate1,
		Map3,
		Map2,
	}

	options := hand.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}

func TestPickableCardWhenFirstCardOnTableIsEscape2(t *testing.T) {
	var table Table
	var hand Hand
	var card Card

	table.cards = []Card{
		card.fromId(Escape1),
		card.fromId(Escape2),
	}

	hand.cards = []Card{
		card.fromId(Chest1),
		card.fromId(Roger5),
		card.fromId(Pirate1),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Chest1,
		Roger5,
		Pirate1,
		Map3,
		Map2,
	}

	options := hand.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}

func TestPickableCardWhenFirstCardOnTableIsEscape3(t *testing.T) {
	var table Table
	var hand Hand
	var card Card

	table.cards = []Card{
		card.fromId(Escape1),
		card.fromId(Escape2),
		card.fromId(Chest1),
	}

	hand.cards = []Card{
		card.fromId(Chest1),
		card.fromId(Roger5),
		card.fromId(Pirate1),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Chest1,
		Pirate1,
	}

	options := hand.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}

func TestPickableCardWhenFirstCardOnTableIsCharacter1(t *testing.T) {
	var table Table
	var hand Hand
	var card Card

	table.cards = []Card{
		card.fromId(Pirate1),
	}

	hand.cards = []Card{
		card.fromId(Chest1),
		card.fromId(Roger5),
		card.fromId(Pirate1),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Chest1,
		Roger5,
		Pirate1,
		Map3,
		Map2,
	}

	options := hand.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}

func TestPickableCardWhenFirstCardOnTableIsCharacter2(t *testing.T) {
	var table Table
	var hand Hand
	var card Card

	table.cards = []Card{
		card.fromId(Pirate1),
		card.fromId(Chest2),
	}

	hand.cards = []Card{
		card.fromId(Chest1),
		card.fromId(Roger5),
		card.fromId(Pirate1),
		card.fromId(Map3),
		card.fromId(Map2),
	}

	expected := []CardId{
		Chest1,
		Roger5,
		Pirate1,
		Map3,
		Map2,
	}

	options := hand.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.")
	}
}
