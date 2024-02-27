package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"skullking/models"
	"time"
)

const GamesTable = "games"

type gameRepository struct {
	db *mongo.Database
}

func NewGameRepository(db *mongo.Database) models.GameRepository {
	return gameRepository{
		db: db,
	}
}

func (ur gameRepository) Create(game *models.Game) error {
	type Player struct {
		Id    string `bson:"id"`
		Score int    `bson:"score"`
		Index int    `bson:"index"`
	}

	type Card struct {
		PlayerId string `bson:"player_id"`
		CardId   uint16 `bson:"card_id"`
	}

	type Bid struct {
		PlayerId string `bson:"player_id"`
		Bid      int    `bson:"bid"`
	}

	type Score struct {
		PlayerId string `bson:"player_id"`
		Score    int    `bson:"score"`
	}

	type Trick struct {
		Number             int    `bson:"number"`
		PickedCards        []Card `bson:"picked_cards"`
		WinnerPlayerId     string `bson:"winner_player_id"`
		WinnerCardId       uint16 `bson:"winner_card_id"`
		StarterPlayerIndex int    `bson:"starter_player_index"`
	}

	type Round struct {
		Number             int     `bson:"number"`
		DealtCards         []Card  `bson:"dealt_cards"`
		Bids               []Bid   `bson:"bids"`
		Tricks             []Trick `bson:"tricks"`
		StarterPlayerIndex int     `bson:"starter_player_index"`
		Scores             []Score `bson:"scores"`
	}

	type Game struct {
		Id        primitive.ObjectID      `bson:"_id"`
		Players   []Player                `bson:"players"`
		Rounds    [len(game.Rounds)]Round `bson:"rounds"`
		CreatorId string                  `bson:"creator_id"`
		CreatedAt primitive.DateTime      `bson:"created_at"`
	}

	var rounds [len(game.Rounds)]Round

	for index, round := range game.Rounds {
		var dealtCards []Card

		for playerId, cardIds := range round.DealtCards {
			for _, cardId := range cardIds {
				dealtCards = append(dealtCards, Card{
					PlayerId: playerId,
					CardId:   uint16(cardId),
				})
			}
		}

		var bids []Bid

		round.Bids.Range(func(playerId string, bid int) bool {
			bids = append(bids, Bid{
				PlayerId: playerId,
				Bid:      bid,
			})
			return true
		})

		var scores []Score

		for playerId, score := range round.Scores {
			scores = append(scores, Score{
				PlayerId: playerId,
				Score:    score,
			})
		}

		var tricks = make([]Trick, len(round.Tricks))

		for index, trick := range round.Tricks {
			var pickedCards []Card

			for _, pickedCard := range trick.PickedCards {
				pickedCards = append(pickedCards, Card{
					PlayerId: pickedCard.PlayerId,
					CardId:   uint16(pickedCard.CardId),
				})
			}

			tricks[index] = Trick{
				Number:             trick.Number,
				PickedCards:        pickedCards,
				WinnerPlayerId:     trick.WinnerPlayerId,
				WinnerCardId:       uint16(trick.WinnerCardId),
				StarterPlayerIndex: trick.StarterPlayerIndex,
			}
		}

		rounds[index] = Round{
			Number:             round.Number,
			DealtCards:         dealtCards,
			Bids:               bids,
			Tricks:             tricks,
			StarterPlayerIndex: round.StarterPlayerIndex,
			Scores:             scores,
		}
	}

	var players []Player

	game.Players.Range(func(_ string, player *models.Player) bool {
		p := Player{
			Id:    player.Id,
			Score: player.Score,
			Index: player.Index,
		}
		players = append(players, p)
		return true
	})

	gameId, err := primitive.ObjectIDFromHex(game.Id)

	if err != nil {
		return err
	}

	createdAt := primitive.NewDateTimeFromTime(
		time.Unix(game.CreatedAt, 0),
	)

	g := Game{
		Id:        gameId,
		Players:   players,
		Rounds:    rounds,
		CreatorId: game.CreatorId,
		CreatedAt: createdAt,
	}

	_, err = ur.db.Collection(GamesTable).InsertOne(context.Background(), g)

	if err != nil {
		return err
	}

	return nil
}
