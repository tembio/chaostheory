
EventFactory is in charge of creating new events (users or bets) with the parameters specified in PossibleBetValues
In this case PossibleBetValues gets its values from the config.

EventGenerator generates events periodically with the method RunEventGeneration that receives a function to execute code for each generated event.
The periodicity and other arguments can be specified when the EventGenerator is created. 
Events will always be generated this way:
- first an event user has to be generated to create at least one user
- events to create bets and win/loss events will always use exisiting users
- win/loss events will always refer to bets that have already been created (so, no invalid win/loss events will be sent)


The interface Sender has a method send, and abstrat the way and the destination where we can send our messages.
In this case the interface has an implementation to send rabbitMQ messages.
There could be a mock for tests, or another one to, for example, print messages in console for debugging.


Run Rabbit:
docker build -f Dockerfile.rabbitmq -t myrabbit .
docker run -p 5672:5673 -p 15672:15673 myrabbit

The management UI will be at http://localhost:15673 (user: guest, password: guest).


Improvements:
- Add logging
TODO:
- handle new competition events and add the new competitios to event_factory