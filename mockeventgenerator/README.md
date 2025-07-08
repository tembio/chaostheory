# Mock Server Generator

## Architecture

**EventFactory** is responsible for creating new events (users or bets) with the parameters specified in `PossibleBetValues`.  
In this case, `PossibleBetValues` gets its values from the config.

**EventGenerator** generates events periodically with the method `RunEventGeneration`, which receives a function to execute code for each generated event.  
The periodicity and other arguments can be specified when the `EventGenerator` is created.

Events will always be generated in this order:
- First, a user event is generated to create at least one user.
- Events to create bets and win/loss events will always use existing users.
- Win/loss events will always refer to bets that have already been created (so, no invalid win/loss events will be sent).

The interface **Sender** has a method `send`, and abstracts the way and the destination where we can send our messages.  
In this case, the interface has an implementation to send RabbitMQ messages.  
There could be a mock for tests, or another one to, for example, print messages in the console for debugging.

---

## Run RabbitMQ

```sh
docker build -f Dockerfile.rabbitmq -t myrabbit .
docker run -p 5672:5673 -p 15672:15673 myrabbit
```

### Improvements and TODOs
- Add logging
- handle new competition events and add the new competitios to event_factory