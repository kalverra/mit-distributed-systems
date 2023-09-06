# Lecture 1: Introduction

## Summary

A lot of these lessons I personally absorbed as a part of work, so not a lot of heavy mental lifting here. I got to read a lot about MapReduce during my time at FINRA.

We established what distributed systems are, the categories of the issues they run into, and introduce a couple ways to solve them.

[Video](https://youtu.be/cQP8WApzIQQ?si=ys60ULiwMk_nvxmE) | Read [MapReduce Paper](https://pdos.csail.mit.edu/6.824/papers/mapreduce.pdf)

## Notes

A distributed system is lots of separate computers working together, and if possible, should not exist. They're complicated, expensive, and hard to program. But sometimes necessary.

### When To Build a Distributed System

* We need high performance, and high levels of parallelism
* Fault tolerance
* Physical separations
* Some sort of security goal

### Why Are They Hard?

* Shit's complicated
* They have lots of complex parts executing concurrently; race conditions are a bitch
* Single computer systems usually have 2 states, working and broken. Dist systems can be broken in a ton of fun and complex ways

Distributed systems used to be an academic curiosity, but now they're everywhere. You want to build something as big as Netflix or Amazon, you gonna need a shit load of computers.

### Infra for Applications

We're workin hard here on creating distributed systems that provide:

* Storage
* Communication
* Computation

With the benefits of distributed systems, but with abstractions that sweep the complexity under the rug. Implementations include:

* RPC
* Threads
* Concurrency control

There's some core concepts and considerations when building these systems.

### Scalability; But Does It Scale?

It better. If we're building this thing, it should be in a way that adding more machines increases performance in a non-logarithmic curve. Ideally linear, but exponential is always something to dream for.

### Fault Tolerance

Shit breaks. If your computer fails at a 1% rate, and you've got a data center with thousands of them, that 1% shows up a lot. Scale turns event the rarest issues into constant problems.

What does it mean to be fault tolerant?

* **Availability**: Even if there are failures, we can keep providing service to users like nothing happened. If a server goes down, no one would know besides the engineers.
* **Recoverability**: If something goes disastrously wrong (full power out), we can pick up right where we left off.

To achieve the above, we'll need lots of hard drives and replication. There's a lot of distributed systems thought around how you can most efficiently write to disk as it can be a huge bottleneck (less with modern SSDs though). You also need to make sure you can have that data be the same in multiple places, getting out of sync is bad news.

### Consistency

If I call GET or PUT, it should actually work as expected. If I `PUT K:V` and then `GET K`, I should get V back. This seems basic, but is a surprisingly tricky issue in big systems. If I have 2 replicas, and in the middle of my `PUT` a power outage happens and now the data is out of sync.

"Consistent" can mean different things to different people. **Strong consistency** usually means what you think it means, my `GET` will always get the latest, most correct value, but that tends to be an expensive thing to guarantee. There's also **weak consistency** where, hey, you might get the latest, you might not, depends on the flavor of application I'm working on.
