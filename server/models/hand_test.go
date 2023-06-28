package models

import (
	"reflect"
	"testing"
)

func TestPickableCardWhenTableIsEmpty(t *testing.T) {
	var table Table
	var hand Hand

	hand.cards = []Card{
		newCardFromId(Parrot2),
		newCardFromId(Roger5),
		newCardFromId(SkullKing),
		newCardFromId(Map3),
		newCardFromId(Map2),
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

	table.cards = []Card{
		newCardFromId(Parrot3),
	}

	hand.cards = []Card{
		newCardFromId(Parrot2),
		newCardFromId(Roger5),
		newCardFromId(SkullKing),
		newCardFromId(Map3),
		newCardFromId(Map2),
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

func TestPickableCardWhenFirstCardOnTableIsSuit2(t *testing.T) {
	var table Table
	var hand Hand

	table.cards = []Card{
		newCardFromId(Map14),
		newCardFromId(Parrot3),
	}

	hand.cards = []Card{
		newCardFromId(Parrot1),
		newCardFromId(Map1),
	}

	expected := []CardId{Map1}

	options := hand.pickables(table)

	if !reflect.DeepEqual(expected, options) {
		t.Error("Wrong pickable cards.", options)
	}
}

func TestUserCanPickAnyCardIfNoCardMatchesTheSuit(t *testing.T) {
	var table Table
	var hand Hand

	table.cards = []Card{
		newCardFromId(Parrot3),
	}

	hand.cards = []Card{
		newCardFromId(Chest1),
		newCardFromId(Roger5),
		newCardFromId(Pirate1),
		newCardFromId(Map3),
		newCardFromId(Map2),
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

	table.cards = []Card{
		newCardFromId(Escape1),
	}

	hand.cards = []Card{
		newCardFromId(Chest1),
		newCardFromId(Roger5),
		newCardFromId(Pirate1),
		newCardFromId(Map3),
		newCardFromId(Map2),
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

	table.cards = []Card{
		newCardFromId(Escape1),
		newCardFromId(Escape2),
	}

	hand.cards = []Card{
		newCardFromId(Chest1),
		newCardFromId(Roger5),
		newCardFromId(Pirate1),
		newCardFromId(Map3),
		newCardFromId(Map2),
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

	table.cards = []Card{
		newCardFromId(Escape1),
		newCardFromId(Escape2),
		newCardFromId(Chest1),
	}

	hand.cards = []Card{
		newCardFromId(Chest1),
		newCardFromId(Roger5),
		newCardFromId(Pirate1),
		newCardFromId(Map3),
		newCardFromId(Map2),
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

	table.cards = []Card{
		newCardFromId(Pirate1),
	}

	hand.cards = []Card{
		newCardFromId(Chest1),
		newCardFromId(Roger5),
		newCardFromId(Pirate1),
		newCardFromId(Map3),
		newCardFromId(Map2),
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

	table.cards = []Card{
		newCardFromId(Pirate1),
		newCardFromId(Chest2),
	}

	hand.cards = []Card{
		newCardFromId(Chest1),
		newCardFromId(Roger5),
		newCardFromId(Pirate1),
		newCardFromId(Map3),
		newCardFromId(Map2),
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
