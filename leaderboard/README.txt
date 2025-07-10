DATA STRCUTURES

type scoresPerCompetition map[uint]usersIDToUser // map[competitionID]map[userID]User

This allows for fast lookups per user when we need to find a user to update their total score


TDOO:
- calculate rate with exchange rate
- make sure logic is ok with example



RUN 
docker build -t leaderboard-service .
docker run --rm leaderboard-service




------- REQUEST TO CREATE COMPETITIONS

curl -X POST http://localhost:8080/competitions \
  -H "Authorization: Bearer secrettoken" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "July Challenge",
    "score_rule": "event_type==\"bet\" ? amount : 0",
    "start_time": "2025-07-01T00:00:00Z",
    "end_time": "2025-07-31T23:59:59Z",
    "rewards": {"1": 100, "2-5": 50, "6+": 20}
  }'



  curl -X POST http://localhost:8080/competitions \
  -H "Authorization: Bearer secrettoken" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Poker Weekend",
    "score_rule": "event_type==\"bet\" && game==\"Poker\" ? amount : 0",
    "start_time": "2025-07-12T00:00:00Z",
    "end_time": "2025-07-13T23:59:59Z",
    "rewards": {"1": 300, "2-3": 100}
  }'



  curl -X POST http://localhost:8080/competitions \
  -H "Content-Type: application/json" \
  -d '{"name":"No Auth"}'



---- REQUESTS TO GET N TOP USERS FROM COMPETITIONS
curl -X GET http://localhost:8080/leaderboards/1
curl -X GET "http://localhost:8080/leaderboards/2?count=5"

invalid id
curl -X GET http://localhost:8080/leaderboards/abc

  TODO::
  logs instead of prints
  not expose internal error (for example DB)
  proper checks for validity of data