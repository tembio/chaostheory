

TDOO:
- calculate rate with exchange rate
- make sure logic is ok with example



RUN 
docker build -t leaderboard-service .
docker run --rm leaderboard-service






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