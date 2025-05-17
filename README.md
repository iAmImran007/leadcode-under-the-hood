# 🧠 LeadCode Under The Hood

A minimal backend system to compile and judge C++ code submissions, inspired by platforms like **LeetCode**. Built with Go and Docker.

I’ve always been curious about how platforms like LeetCode compile and run user-submitted code in different languages securely. This project is my attempt to explore and build a simple, backend-focused version of an online judge system to understand what happens **under the hood**.

---

## 🚀 How It Works

This is a simplified version of an online judge system:

1. Problems and associated test cases are stored in a database.
2. There are three main API endpoints:
   - `GET /problems`: Returns all problems.
   - `GET /problems/{id}`: Returns details of a specific problem.
   - `POST /submit`: Accepts user code, spins up a Docker container, **compiles and executes the code securely inside the container**, and compares output with expected results.
3. Inside the `/submit` endpoint:
   - The submitted code is written to a temporary file.
   - A Docker container compiles the C++ code using `g++`.
   - The compiled binary is executed against each test case with a timeout.
   - Execution results are compared to expected outputs.
   - A summary result (pass/fail) is returned.

---

## 🛠 Features

- ✅ C++ code submission support
- 🐳 Test case validation inside isolated Docker containers
- 🧩 Auto-populated database with sample problems at startup
- 🧠 Minimalistic code judge logic (no frontend/UI yet)

---

## 📦 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/iAmImran007/leadcode-under-the-hood.git
cd leadcode-under-the-hood
