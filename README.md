# Reddit Clone with Actor Model in Go

A Reddit-inspired forum application built in Go, leveraging the Actor Model (ProtoActor) for concurrency, scalability, and fault-tolerance.
This project demonstrates how distributed actor-based systems can be applied to real-world applications like Reddit-style communities.

# Features

User registration & authentication

Post creation and threaded discussions

Commenting & voting system

Actor model concurrency for handling user sessions, posts, and votes

Modular folder structure (core, handlers, schemas, templates, tests)

Easy to extend and scale

Tech Stack

Language: Go

Concurrency Model: Actor Model (ProtoActor)

Templates: HTML templates for UI

Schemas: JSON/Protobuf schemas for structured communication

Testing: Unit tests under /tests

# Project Structure
.
├── core/          # Core business logic & actor definitions
├── handlers/      # HTTP handlers (routes & controllers)
├── reddit-clone/  # Reddit specific features
├── schemas/       # Protobuf/JSON schemas
├── templates/     # HTML templates for frontend
├── tests/         # Unit and integration tests
├── go.mod         # Go module definition
├── go.sum         # Dependency checksum
├── LICENSE        # License information
├── main.go        # Entry point of the application
└── README.md      # Project documentation

# Installation & Setup

Clone the repository:

git clone https://github.com/<your-username>/<repo-name>.git
cd <repo-name>


# Install dependencies:

go mod tidy


Run the application:

go run main.go

Why Actor Model?

The Actor Model provides a simple yet powerful abstraction for building concurrent and distributed systems.
In this project:

Each post, user, and vote can be modeled as an actor.

Actors communicate via messages, avoiding shared memory issues.

The system becomes more scalable and fault-tolerant.

Future Improvements

Add persistent storage (Postgres/MongoDB)

Implement upvote/downvote ranking system

Enhance UI with React/Vue frontend

Deploy on Docker/Kubernetes

# License

This project is licensed under the MIT License.
