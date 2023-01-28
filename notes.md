# Entities

## Game
- id: string
- state: string
  - start
  - bidding
  - picking
  - calculating
  - end
- round: int
- clients: map[userId]*Client
- pickingUserId: int

## User
- name: string
- email: string
- email_verified_at: datetime
- game_id: int (this will be used to reconnect user after login)
- created_at

## Play
- round
- score
- user_id
- game_id
- card_id