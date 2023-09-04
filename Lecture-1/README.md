# Lecture 1: Introduction

[Video](https://youtu.be/cQP8WApzIQQ?si=ys60ULiwMk_nvxmE) | Read [MapReduce Paper](https://pdos.csail.mit.edu/6.824/papers/mapreduce.pdf)

A distributed system is lots of separate computers working together, and if possible, should not exist. They're complicated, expensive, and hard to program. But sometimes necessary.

## When To Build a Distributed System

* We need high performance, and high levels of parallelism
* Fault tolerance
* Physical separations
* Some sort of security goal

## Why Are They Hard?

* Shit's complicated
* They have lots of complex parts executing concurrently; race conditions are a bitch
* Single computer systems usually have 2 states, working and broken. Dist systems can be broken in a ton of fun and complex ways

Distributed systems used to be an academic curiosity, but now they're everywhere. You want to build something as big as Netflix or Amazon, you gonna need a shit load of computers.
