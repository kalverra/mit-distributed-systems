# Map Reduce

Doing it with MIT's actual example project would probably be better, but I so very much deeply despise those types of academic "your code here" projects for reasons I can't sufficiently articulate. So fuck it, do my own then.

Each worker is going to pretend to be its own machine, communicating using RPC with the coordinator. The lab has a couple happy paths to test for correctness, and has a few fault scenarios to deal with.

- [ ] Run `Map` and `Reduce` all in parallel
- [ ] Handle worker code crashes
- [ ] Handle long-running workers
- [ ] Bonus: Handle coordinator crashes

> The coordinator should notice if a worker hasn't completed its task in a reasonable amount of time (for this lab, use ten seconds), and give the same task to a different worker.
