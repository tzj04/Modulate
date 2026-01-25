# Modulate

**Author:** Tan Zheng Jia
**Application Type:** Full-Stack Discussion Platform

---

## 🚀 Overview

Modulate is a modular forum application designed for structured discussions. It features a high-performance **Go** backend and a type-safe **React (TypeScript)** frontend. The application emphasizes data integrity through "Soft Deletion," ensuring that threaded conversations remain readable even if a user removes their specific contribution.

### Key Features

- **Modular Organization:** Posts are grouped by specific modules (e.g., "General," "Development," "Feedback").
- **Threaded Discussions:** Nested, recursive comment sections allowing for deep conversations.
- **Soft Delete Logic:** Posts and comments are "ghosted" (content scrubbed) rather than physically removed to preserve thread context and database integrity.
- **Backend Authorization:** All write operations (Update/Delete) verify the user's identity via JWT against the resource owner in the database.

---

## 🛠️ Tech Stack

- **Frontend:** React, TypeScript
- **Backend:** Go
- **Database:** PostgreSQL

---

## ⚙️ Setup & Installation

### 0. Enviroment

1.  copy the .env.example files and fill in the required details
    ```bash
    cp .env.example .env
    cp ./backend/.env.example ./backend.env
    cp ./frontend/.env.example ./frontend/.env
    ```

### 1. Database Setup

Ensure you have **PostgreSQL** installed and running.

1.  Create a database named `modulate_db`:
    ```sql
    CREATE DATABASE modulate_db;
    ```
2.  Execute the schema migrations in backend/internal/db/migrations/init.sql

3.  Go to nusmods_data_processing and run:
    ```python
    python extract.py
    python transform.py
    python load.py
    ```

### 2. Backend Setup

1.  Navigate to the backend directory:
    ```bash
    cd backend
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```
3.  **Environment Configuration:** Ensure your database connection string in .env matches your local Postgres credentials:

4.  Run the application:
    ```bash
    go run main.go
    ```
    _The server will start on `http://localhost:8080` by default unless specified._

### 3. Frontend Setup

1.  Navigate to the frontend directory:
    ```bash
    cd frontend
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Start the development server:
    ```bash
    npm run dev
    ```
    _The frontend will be accessible at `http://localhost:3000` by default unless specified._

---

## 📖 Grading / Evaluation Guide

To test the core logic of the application, follow these steps:

1.  **Auth Flow:** Register a new user and log in.
2.  **Creation:** Navigate to a Module and create a new Post.
3.  **Threading:** Add a comment to your post, then add a reply to that comment to verify recursive rendering.
4.  **Soft Delete:** Delete your post or comment.
    - **Observation:** The post should remain accesible via the link but show `[Deleted]` as the title and `[This post has been removed by the author]` as the content.
5.  **Security Verification:** \* Log out and try to access a protected route (e.g., `/api/posts/{id}`).
    - Log in as a different user and attempt to edit/delete a post you don't own. The backend will return a `403 Forbidden` because it verifies the `userID` from the JWT context.

---

## AI Declaration

1. Some of the styling in the css was done by GPT to speed up the process
2. When there were parts I wasn't too sure about, I clarified with GPT to make sure my logic was sound
