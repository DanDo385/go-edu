# Go Systems Programming Course — Teaching & Structure Instructions

This repository is a **CS50-style Go systems programming course**, taught through real, production-quality code.

The goal is for learners to understand **how Go works from first principles**, not just how to write syntactically correct programs.

---

## Course Philosophy

* Teach **systems concepts first**, Go syntax second
* Favor **mental models** over rules of thumb
* Use **real codebases**, not contrived examples
* Progress from **small, isolated concepts** → **integrated systems**

Learners should finish this course able to:

* Read unfamiliar Go code with confidence
* Reason about memory, values, references, and concurrency
* Understand *why* Go behaves the way it does

---

## Course Structure

The course is divided into two complementary tracks:

### `minis/` — Focused Concept Labs

* Small, self-contained programs
* Each mini introduces **one primary concept**
* Minimal code, maximum clarity
* Used to build intuition before scaling up

Examples:

* Values vs pointers
* Structs and method receivers
* Interfaces as behavior contracts
* Goroutines and channels in isolation

---

### `geth/` — Integrated Systems Track

* Larger, realistic systems
* Concepts from `minis` are **reused, extended, and composed**
* Emphasis on orchestration, lifecycle, and data flow

Examples:

* Long-running processes
* Concurrency coordination
* State management over time
* Interaction between components

---

## Lesson Design Rules

Each lesson must:

1. Introduce **one or two new ideas only**
2. Explicitly connect to **previous lessons**
3. Explain *why this abstraction exists*
4. Avoid overwhelming the learner with comments or boilerplate

Lessons should **build, not branch**.

---

## Commentary & Documentation Guidelines

* File-level comments explain **system role**
* Section-level comments explain **flow and sequencing**
* Function-level comments explain **assumptions and invariants**
* Inline comments are used **sparingly** and only when clarity demands it

The code should remain **visually readable without scrolling past walls of comments**.

---

## Pointer & Reference Emphasis (Critical)

Whenever pointers or references appear:

* Slow down
* Explain what lives in memory
* Explain what is copied vs shared
* Explain why Go made this tradeoff

If the learner misunderstands pointers, the lesson has failed.

---

## Outcome

This course should feel like:

**“CS50, but for Go systems programming—using real-world code.”**
