# Grader: An Automated Code Grader System

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/your_username/Grader)
![GitHub last commit](https://img.shields.io/github/last-commit/your_username/Grader)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/your_username/Grader)
![Docker Pulls](https://img.shields.io/docker/pulls/your_username/Grader)

## Table of Contents

- [Introduction](#introduction)
- [Technology Stack](#technology-stack)
- [Service Architecture](#service-architecture)
- [Getting Started](#getting-started)
- [Contribute](#contribute)

## Introduction

Grader is a versatile, automated code grader system designed to streamline the process of homework verification. It is perfect for educational institutions, coding bootcamps, or any organization looking to simplify the grading process of coding assignments. Grader is fully capable of receiving user-submitted solutions, adding test files for further homework verification, and more.

## Technology Stack

Grader leverages an array of modern technologies to provide robust and scalable service:

- **Go (Golang)**: A statically typed, compiled language with syntax similar to C but with garbage collection, memory safety features, and CSP-style concurrency.
- **Go-chi**: A lightweight, idiomatic, and composable router for building Go HTTP services.
- **PostgreSQL**: An advanced, enterprise-class, and open-source relational database system.
- **Redis**: An open-source, in-memory data structure store used as a database, cache, and message broker.
- **RabbitMQ**: An open-source message broker software that implements the Advanced Message Queuing Protocol (AMQP).
- **Docker**: A platform to develop, ship, and run applications segregated into separate container instances.

## Service Architecture

Grader comprises of three key services:

1. **Grader Service**: This is where solutions are received, validated, and processed. The solution file is transferred via Docker volume mount to a container, where the task is evaluated.
2. **Queue Service**: Manages the distribution and orchestration of tasks using RabbitMQ as the underlying message broker.
3. **Server (User Part)**: Handles user interactions, including logins, solution uploads, and accessing test files.

## Getting Started

```sh
# Clone the repository
git clone https://github.com/ALIBRAGIMOV/Grader.git

# Go to the project directory
cd Grader
