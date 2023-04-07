# Entities

## Game

```
{
  "id": int,
  "round": int,
  "trick": int,
  "state": enum(start, bidding, picking, calculating, ended),
  "expiration_time": int,
  "users": [
    {
      "id": int,
    }
  ],
  "rounds": [
    {
      "number": int,
      "scores": [
        {
          "user_id": int,
          "score": int
        }
      ],
      "dealt_cards": [
        {
          "user_id": int,
          "card_ids": []int
        }
      ],
      "bids": [
        {
          "user_id": int,
          "value": int
        }
      ],
      "tricks": [
        {
          "number": int,
          "picking_user_id": int,
          "picked_cards": [
            {
              "user_id": int,
              "card_id": int
            }
          ]
        }
      ]
    }
  ]
}
```

## User
```
{
  name: string,
  email: string,
  game_id: int,
  verified_at: datetime,
  crated_at: datetime
}
```