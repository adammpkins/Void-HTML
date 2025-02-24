# Void-HTML

**Void-HTML** is a lightweight microblogging and social networking platform built with Go, Fiber, and server-rendered HTML. Inspired by the "Laravel Bootcamp" approach, Void-HTML takes a deep dive into building robust web applications using Go's performance and simplicity while adhering to clean architecture principles.

This repository is the HTML version of the Void platform. A separate repository (e.g., Void-Inertia or Void-Vue) is planned for a more dynamic, client-driven front-end using Inertia/Vue.

---

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Architecture & Design](#architecture--design)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Running the Application](#running-the-application)
- [Project Structure](#project-structure)
- [Middleware & Session Management](#middleware--session-management)
- [Notifications & Event-Driven Architecture](#notifications--event-driven-architecture)
- [Development Workflow](#development-workflow)
- [Contributing](#contributing)
- [License](#license)
- [Additional Notes](#additional-notes)

---

## Overview

Void-HTML is designed to be a minimalistic yet fully functional microblogging platform. Its primary focus is to:

- Allow users to register, login, and manage their profiles.
- Enable users to post "shouts" (short messages) and view a global feed.
- Provide a robust notification system using RabbitMQ for real-time interactions.
- Serve as a learning ground for modern web development patterns, inspired by frameworks like Laravel but leveraging the speed of Go.

---

## Key Features

- **User Authentication & Session Management:**  
  Secure registration, login, and logout functionality. Sessions are managed via custom middleware to ensure that user credentials and sessions are managed reliably.

- **Server-Rendered HTML:**  
  Uses Fiber's HTML templating engine to render pages on the server side, making it a great starting point for projects that rely on traditional multi-page application (MPA) patterns.

- **Real-Time Notifications:**  
  Leverages RabbitMQ to implement an asynchronous notification system that informs users about new shouts and interactions.

- **Microblogging Capabilities:**  
  Users can post shouts, view a global feed, and interact through "echoes" (replies), all managed via RESTful endpoints and rendered templates.

- **Profile Customization:**  
  Each user has a customizable profile complete with a bio and avatar upload functionality. Images are processed using the imaging library for resizing.

- **ORM-Based Data Management:**  
  Built on GORM and SQLite for efficient data storage and retrieval with minimal configuration.

---

## Architecture & Design

### Modular Structure

Void-HTML is structured into clearly defined modules:
- **cmd/server:** Contains the main application entry point.
- **internal:** Houses the core business logic including:
  - **db:** Database initialization and migrations.
  - **events:** Event definitions and publisher.
  - **handlers:** HTTP route handlers for different functionalities (authentication, user management, shouts, notifications).
  - **middleware:** Custom middleware for session management and authentication checks.
  - **models:** GORM models for Users, Shouts, Echoes, and Notifications.
  - **services/notifications:** Logic for processing and sending notifications.
- **pkg:** Contains reusable packages like RabbitMQ integration and session management.
- **web:** Contains static assets and HTML templates.

### Middleware & Session Management

User authentication is managed using middleware:
- **GetUserFromSession:** Retrieves the user ID from the session (using a custom session package) and stores it in the request context under `c.Locals("UserID")`.
- **RequireLogin:** Ensures that endpoints requiring authentication have a valid user ID. If not, it redirects to the login page.

This flow ensures that all route handlers can safely rely on the presence of a valid user ID without directly handling session logic.

### Notifications & Event-Driven Architecture

When a new shout is created:
1. The shout is saved to the database.
2. An event is generated (conforming to the `ShoutEvent` interface).
3. The event is marshaled to JSON and published via RabbitMQ.
4. A separate service consumes these notifications and sends real-time alerts to users (excluding the shout's author).

This decoupled, event-driven architecture allows for scalability and a responsive user experience.

---

## Technologies Used

- **Go:** The primary programming language, offering simplicity and performance.
- **Fiber:** A fast, Express-inspired web framework for Go.
- **GORM:** An ORM library that simplifies database interactions.
- **SQLite:** A lightweight and self-contained SQL database engine.
- **RabbitMQ:** A robust message broker used for asynchronous event handling and notifications.
- **HTML Templating:** Server-side rendering using HTML templates for dynamic content delivery.
- **Imaging:** A Go package for image processing (resizing avatars).

---

## Getting Started

### Prerequisites

- **Go 1.23.1** or later installed.
- **RabbitMQ:** Running locally or accessible via network.
- Basic familiarity with Go modules and command-line operations.

### Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/yourusername/Void-HTML.git
   cd Void-HTML
   ```

2. **Download dependencies:**

   ```sh
   go mod download
   ```

3. **Configuration:**

   Update the RabbitMQ connection string in `cmd/server/main.go` if needed:

   ```go
   if err := rabbitmq.Init("amqp://guest:guest@localhost:5672/"); err != nil { ... }
   ```

   The application uses SQLite, and the database file (`void.db`) will be created in the root directory.

### Running the Application

You can run the application using Air for live reloading during development:

```sh
air
```

Or build and run manually:

```sh
go build -o void ./cmd/server
./void
```

Then open your browser and navigate to `http://localhost:3000`.

---

## Development Workflow

- **Live Reload with Air:**  
  Air monitors changes in your Go files, HTML templates, and static assets, automatically rebuilding and restarting the server.

- **Testing & Debugging:**  
  Use logging (Fiber's built-in logging and custom log messages) to trace requests and debug issues in middleware and handlers.

- **Code Organization:**  
  Follow the established modular structure to ensure that your business logic, data models, and routing remain decoupled and maintainable.

---

## Contributing

Contributions are welcome! To get started:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Ensure your changes follow the project's coding standards and pass all tests.
4. Open a pull request with a detailed description of your changes.

For major changes, please open an issue first to discuss what you would like to change.

---

## License

This project is licensed under the MIT License. See the LICENSE file for details.

---

## Additional Notes

- **Inspired by Laravel Bootcamp:**  
  The project draws inspiration from the Laravel Bootcamp series, translating many of its best practices and architectural decisions into the Go ecosystem.

- **Future Plans:**  
  A version using Inertia/Vue is planned, which will provide a more dynamic, SPA-like front end while keeping a unified backend API.

- **Feedback & Support:**  
  Feel free to open issues if you encounter any bugs or have suggestions for improvements. Your feedback is valuable!
