# Structure
The project has 3 services:
- mockeventgenerator: generates the events and send them to the rabbitMQ queue.
- leaderboard: reads the events from the queue and calculates the results of the leaderboards. 
                It also exposes two HTTP enpodints, one to create competitions, and another one to create users.
- frontend: is a webserver to display the competitions. It gets updated in realtime when a new event is received.

The shared code between services is in a separate go module in the `common` folder.

# Building and running the project

I've provided two ways of building/running the project:
    - Locally, using docker compose. The services will run in docker images and a local network will be created so that
      they can communicate.
      In order to build the images and run the project we have to run the `build_compose.sh` (giving execution permissions first).
      The webpage for the leaderboard will be in http://localhost:8081/
      The RabbitMQ console can be accessed in http://localhost:15672/

      If we kill the leaderboard service, it will NOT start automatically, for that we need to use the other way running the project using Kubernetes explained below. 
      If we restart the leaderboard manually it will keep the state, and will continue working normally.

    - The other way, is using kubernetes, having a local cluster with minikube.
      In order to do that,we need to run the `build_k8s.sh` script. That will build the docker images and create the necessary
      deployments and services to maintain always the `leaderboard` service running, recreating it if it is shutdown manually.
      As mentioned before, having minikube is a requirement for this setup, it needs to be installed and running:
      `minikube start --cpus 8 --memory 10g`
      
      (RabbitMQ requires ate least that memory, it fails with less)
    
      In order to be able to access the webpage we need to get the IP and port of the cluster, that can be obtained with the command:

      `minikube service rabbitmq --url` will give access to the rabbit console in the kuberneted client


# Considerations
- The mockeventgenerator parameters can be adjusted in the config file. The requirements of the assigment was that the average number of messages had to be 10/s, but it is now set to 5, becuase it makes it easiet to visualise in the UI.
This can be changed to 10 at any time, along with the values that the events will have.

- If we run the project and make requests to create a new leader board, if the event match the rules provided in the leaderboard
the UI will be uptated and create a new leaderboard. Some example requests for this are the following:

## Create new competitions

```
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
  ```

```
  curl -X POST http://localhost:8080/competitions \
  -H "Authorization: Bearer secrettoken" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Blackjack Weekend",
    "score_rule": "event_type==\"bet\" && game==\"Blackjack\" ? amount : 0",
    "start_time": "2025-07-12T00:00:00Z",
    "end_time": "2025-07-13T23:59:59Z",
    "rewards": {"1": 300, "2-3": 100}
  }'
```

## Request N top users from competitions

```
curl -X GET http://localhost:8080/leaderboards/1
```
```
curl -X GET "http://localhost:8080/leaderboards/2?count=5"
```

# Improvements and TODOs

- I've only used prints instead of a proper logging library
- I'm sending new user events but the leaderboard is not doing anything with them (it was clarified in the requirements later that we didn't need to store them)
- Some internal errors (like, rabbit or DB errors) are being exposed to the API, they should be hidden, that could be improved
- The validation of fields in requests is very basic
- The rules that are compiled for to calculate the machtes in the event can be cached so avoid recompiling

