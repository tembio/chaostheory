# Leaderboard Service

## Overview

The `leaderboard` service is designed to process betting and user events, maintain leaderboards, and expose endpoints for competition management and leaderboard queries. This service is intended to be modular and extensible for future enhancements.

## File Structure

- **main.go**
  - The entry point for the leaderboard service. It is responsible for initializing the service, starting any HTTP servers, and wiring up dependencies. In a full implementation, this file would set up routing, connect to databases or message queues, and launch the leaderboard logic.


## Core Data Types and Components

- **event_receiver**
  - Responsible for receiving and processing incoming events (such as bet, win, loss, and user events) from a message queue. 
    It receives a callbacl with the body of the message and an ack function to achnowledge the message was received.

- **leaderboard**
  - The main data structure that maintains the scores and rankings for users in various competitions. 
  It supports the following operations: updating scores, registering new competitions and loading leaderboards (used to initialise the data 
  from DB)

- **rule_evaluator**
  - Encapsulates the logic for evaluating competition rules. 
    When a competition is registered, its rule is added to the ruleEvaluator list of rules.
    When an event is received the event is evaluated against these rules, and a list of matches is returned.

- **repositories/**
  - Contains data access logic and abstractions for persistent storage. These include:
    - **competition_repository.go**: Manages CRUD operations for competitions.
    - **user_repository.go**: Manages user data.
    - **leaderboard_repository.go**: Handles storage and retrieval of leaderboard data.
