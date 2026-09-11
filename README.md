# 🚀 Event-Driven Log Aggregator in Go

A high-throughput log aggregation system built with **Go**, **NATS**, and **MongoDB**. Designed to capture, buffer, and process application logs asynchronously with low latency and efficient resource management.

---

## 📌 Overview

Modern distributed systems generate high volumes of logs that require centralized collection and durable storage. This project provides a lightweight and scalable backend service that consumes log events from a message queue, processes them concurrently using a worker pool pattern, and persists them into MongoDB.

---

## ✨ Highlights

* **Asynchronous Ingestion:** Uses NATS as an event broker for decoupled message handling.
* **Efficient In-Memory Processing:** Implements a concurrent Worker Pool to batch database writes and optimize resource usage.
* **Clean Architecture:** Modular structure dividing configuration, data models, queue handlers, and database persistence layers.
* **Graceful Shutdown:** Ensures buffered in-memory logs are safely flushed to storage before process termination.
* **Local Development Setup:** Includes Docker Compose for running NATS, MongoDB, and Mongo Express out of the box.

---

## 📂 Project Layout

* `cmd/consumer`: Main service entrypoint for log processing.
* `cmd/producer`: Test utility to generate synthetic traffic.
* `internal/`: Core logic including configuration, models, queue subscribers, worker pools, and repository bindings.

---

## 🛠️ Tech Stack

* **Go**
* **NATS Messaging**
* **MongoDB**
* **Docker / Docker Compose**

---

## 🚀 Quick Start

1. **Start Services:**
   ```bash
   docker-compose up -d