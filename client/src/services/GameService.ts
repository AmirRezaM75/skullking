import type {
	Player,
	Card,
	DealResponse,
	StartBiddingResponse,
	StartPickingResponse,
	BadeResponse,
	EndBiddingResponse,
	PickedResponse,
	AnnounceTrickWinnerResponse,
	NextTrickResponse,
	ReportErrorResponse,
	AnnounceScoresResponse,
	Table,
	InitResponse,
	PlayerResponse,
	LeftResponse
} from './../types';
import { GameCommand, GameState } from './../constants';
import type CardService from './CardService';
import Time from '../utils/Time';
import type Swiper from 'swiper';

class GameService {
	// Authenticated user id
	authId: string;

	cardService: CardService;

	players: Player[];

	// Dealt cards for authenticated user
	cards: Card[] = [];

	table: Table = {
		cards: [],
		hasWinner: false
	};

	state: GameState = GameState.Pending;

	round = 0;

	trick = 1;

	timer = 0;

	notifierMessage = '';

	bids: number[] = [];

	bid = 0;

	// Epoch time when we are not going to process to next commands from server
	// because we are busy to complete CSS animations.
	waiter = 0;

	showCountdown = false;

	countdownColor: 'blue' | 'red' = 'blue';

	creator: {
		id: string;
		username: string;
	} = {
		id: '',
		username: ''
	};

	constructor(cardService: CardService, authId: string) {
		this.players = [];
		this.cardService = cardService;
		this.authId = authId;
	}

	async handle(command: GameCommand, content: any) {
		const now = Date.now();

		if (this.waiter != 0 && this.waiter > now) {
			const time = new Time();
			await time.sleep(this.waiter - now);
			this.waiter = 0;
		}

		if (GameCommand.Init == command) {
			this.init(content);
		}

		if (GameCommand.Joined == command) {
			this.joined(content);
		}

		if (GameCommand.Left == command) {
			this.left(content);
		}

		if (GameCommand.AnnounceScores == command) {
			this.announceScores(content);
		}

		if (GameCommand.Deal == command) {
			this.deal(content);
		}

		if (GameCommand.StartBidding === command) {
			this.startBidding(content);
		}

		if (GameCommand.EndBidding === command) {
			this.endBidding(content);
		}

		if (GameCommand.Bade === command) {
			this.bade(content);
		}

		if (GameCommand.StartPicking === command) {
			this.startPicking(content);
		}

		if (GameCommand.Picked === command) {
			this.picked(content);
		}

		if (GameCommand.AnnounceTrickWinner === command) {
			this.announceTrickWinner(content);
		}

		if (GameCommand.NextTrick === command) {
			this.nextTrick(content);
		}

		if (GameCommand.EndGame == command) {
			this.endGame();
		}

		if (GameCommand.ReportError == command) {
			this.reportError(content);
		}

		return this;
	}

	// we need another function which is triggered after main game handler function.
	// because we need to make sure DOM is updated before working with swiper
	public postHandler(command: GameCommand, deckSwiper: Swiper, tableSwiper: Swiper) {
		if (GameCommand.Picked === command) {
			// Picked card animation takes one second to be complete.
			setTimeout(() => {
				deckSwiper.init();
				deckSwiper.update();
				tableSwiper.init();
				tableSwiper.update();
				tableSwiper.slideTo(this.table.cards.length - 1);
			}, 500);
			const audio = new Audio('/sounds/picked.mp3');
			audio.play();
		}

		if (GameCommand.Init === command) {
			setTimeout(() => {
				// First we need to check if is there any cards in deck and table to avoid following error
				// swiper.el.querySelectorAll is not a function or its return value is not iterable
				if (this.table.cards.length !== 0) {
					tableSwiper.init();
					tableSwiper.update();
				}

				if (this.cards.length !== 0) {
					deckSwiper.init();
					deckSwiper.update();
				}
			}, 500);
		}

		if (GameCommand.AnnounceTrickWinner === command) {
			const winnerCardIndex = this.table.cards.findIndex((card) => card.isWinner);

			if (winnerCardIndex !== -1) {
				tableSwiper.slideTo(winnerCardIndex);
				const audio = new Audio('/sounds/announceTrickWinner.mp3');
				audio.volume = 0.2;
				audio.play();
			}
		}

		if (GameCommand.Deal === command) {
			deckSwiper.init();
		}
	}

	init(content: InitResponse) {
		this.round = content.round;
		this.trick = content.trick;
		this.state = content.state;

		if (content.tableCards) {
			content.tableCards.forEach((tableCard) => {
				const card = this.cardService.findById(tableCard.cardId);
				const player = content.players.find((player) => player.id === tableCard.playerId);

				if (card) {
					card.ownerUsername = player ? player.username : ''
					this.table.cards.push(card);
					// TODO: announceTrickWinner has not been implemented
				}
			});
		}

		content.players.forEach((player) => {
			this.addPlayer(player);

			if (player.handCardIds && player.id == this.authId) {
				player.handCardIds.forEach((cardId) => {
					const card = this.cardService.findById(cardId);
					if (card) {
						this.cards.push(card);
					}
				});
			}
		});

		if (content.state === GameState.Bidding) {
			this.startBidding({
				endsAt: content.expirationTime,
				state: content.state
			});

			const player = content.players.find((player) => player.id === this.authId);
			if (player) {
				this.bid = player.bid;
			}
		}

		if (content.state === GameState.Picking) {
			const player = content.players.find((player) => player.id === this.authId);

			this.startPicking({
				state: content.state,
				endsAt: content.expirationTime,
				playerId: content.pickingUserId,
				cardIds: player ? player.pickableCardIds ?? [] : []
			});
		}

		const creator = this.findPlayerById(content.creatorId);

		if (creator) {
			this.creator = {
				id: creator.id,
				username: creator.username
			};
		}
	}

	joined(content: PlayerResponse) {
		this.addPlayer(content);
	}

	left(content: LeftResponse) {
		this.deletePlayerById(content.playerId);
	}

	// To determine when next round is started, use deal()
	nextTrick(content: NextTrickResponse) {
		this.table.cards = [];
		this.table.hasWinner = false;
		this.round = content.round;
		this.trick = content.trick;
	}

	endGame() {
		this.cards = [];
		this.table.cards = [];
		const winner = this.players.reduce((previous, current) => {
			return previous.score > current.score ? previous : current;
		});

		this.notifierMessage = `${winner.username} won the game`;
	}

	announceScores(content: AnnounceScoresResponse) {
		const audio = new Audio('/sounds/announceScores.mp3');
		audio.play();

		this.players.forEach((player) => {
			content.scores.forEach((item) => {
				if (item.playerId === player.id) {
					player.score = item.score;
				}
			});
		});
	}

	deal(content: DealResponse) {
		this.table.cards = [];
		this.table.hasWinner = false;
		this.state = content.state;
		this.round = content.round;
		this.trick = content.trick;
		// Using map() raises typescript warning because of nullable findById
		content.cards.sort((a,b) => a-b).forEach((cardId) => {
			const card = this.cardService.findById(cardId);
			if (card) {
				this.cards.push(card);
			}
		});

		this.players.forEach((player) => {
			player.wonTricksCount = 0;
			player.bid = 0;
		});
	}

	startBidding(content: StartBiddingResponse) {
		// As the server will bid number 0 by default,
		// It is reasonable to indicate that to player.
		this.bid = 0;
		const now = new Date().getTime() / 1000;
		this.timer = Math.floor(content.endsAt - now);
		this.state = content.state;
		this.bids = [...Array(this.round + 1).keys()];
		this.notifierMessage = 'Bidding Time';
		this.countdownColor = 'blue';
		this.showCountdown = true;
	}

	endBidding(content: EndBiddingResponse) {
		this.showCountdown = false;
		this.bids = [];
		this.state = GameState.EndBidding;
		content.bids.forEach((bid) => {
			const player = this.findPlayerById(bid.playerId)
			if (player) {
				player.bid = bid.number
			}
		});
	}

	bade(content: BadeResponse) {
		this.bid = content.number;
	}

	startPicking(content: StartPickingResponse) {
		this.bids = [];

		const now = new Date().getTime() / 1000;
		this.timer = Math.floor(content.endsAt - now);
		this.state = content.state;
		this.players.forEach((player) => {
			player.picking = false;
			if (player.id === content.playerId) {
				player.picking = true;
				this.notifierMessage =
					player.id === this.authId ? 'Pick your card' : `${player.username} is picking`;
				this.countdownColor = player.id === this.authId ? 'blue' : 'red';

				if (player.id === this.authId) {
					const audio = new Audio('/sounds/startPicking.mp3');
					audio.play();

					this.cards.forEach((card) => {
						if (!content.cardIds.includes(card.id)) {
							card.disabled = true;
						}
					});
				}
			}
		});
		this.showCountdown = true;
	}

	picked(content: PickedResponse) {
		this.players.forEach((player) => {
			player.picking = false;
		});
		this.notifierMessage = '';

		this.cards.forEach((card) => {
			card.disabled = false;
		});

		this.showCountdown = false;
		const index = this.cards.findIndex((card) => card.id === content.cardId);

		if (index !== -1) {
			this.cards.splice(index, 1);
		}

		const card = this.cardService.findById(content.cardId);
		const player = this.findPlayerById(content.playerId);

		if (card) {
			card.ownerUsername = player ? player.username : ''
			this.table.cards.push(card);
		}

		// When this is the last person who picked the card
		// before announcing the trick winner, we need time
		// to make sure picked-card-animation is done
		// before winner animation
		if (this.table.cards.length === this.players.length) {
			const time = new Time();
			this.waiter = time.add(2);
		}
	}

	announceTrickWinner(content: AnnounceTrickWinnerResponse) {
		// We need to reset is_winner which is set from previous announcing
		this.table.cards.forEach((card) => {
			card.isWinner = false;
		});

		this.table.hasWinner = true;

		if (content.playerId === '') {
			this.notifierMessage = `No one won the trick.`;
		} else {
			this.table.cards.forEach((card) => {
				if (card.id === content.cardId) {
					card.isWinner = true;
				}
			});

			this.players.forEach((player) => {
				if (player.id === content.playerId) {
					player.wonTricksCount++;
					this.notifierMessage = `${player.username} Won the trick.`;
				}
			});
		}

		const time = new Time();
		this.waiter = time.add(2);
	}

	reportError(content: ReportErrorResponse) {
		alert(content.message);
	}

	addPlayer(player: PlayerResponse) {
		// When user joins a game, it receives "INIT" and "JOINED" commands
		// at the same time, to avoid inserting auth user data twice
		// we need to check if user id is already exists or exclude
		// authenticated user id but I think using AuthService is overhead in this class
		const exists = this.findPlayerById(player.id);

		if (exists) return;

		const p: Player = {
			avatar: '/images/avatars/' + player.avatar,
			id: player.id,
			username: player.username,
			picking: false,
			bid: player.bid,
			score: player.score,
			wonTricksCount: player.wonTricksCount
		};

		this.players.push(p);
	}

	deletePlayerById(id: string) {
		const index = this.players.findIndex((player) => player.id === id);
		if (index !== -1) {
			this.players.splice(index, 1);
		}
	}

	findPlayerById(playerId: string): Player | null {
		for (let i = 0; i < this.players.length; i++) {
			if (this.players[i].id === playerId) {
				return this.players[i];
			}
		}
		return null;
	}

	findPickingPlayerId(): string {
		if (this.state !== GameState.Picking) {
			return '';
		}

		for (let i = 0; i < this.players.length; i++) {
			if (this.players[i].picking) {
				return this.players[i].id;
			}
		}
		return '';
	}
}

export default GameService;
