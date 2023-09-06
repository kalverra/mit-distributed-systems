# MapReduce Notes and Summary

[MapReduce Paper](https://pdos.csail.mit.edu/6.824/papers/mapreduce.pdf)

Google was dealing with tons of data, and they were effectively trying to sort the entire internet. So they bought all the computers, and had engineers ad-hoc design programs to make use of them for each problem. This was tricky, and required some framework to utilize these computers without considering too much detailed complexity.

A MapReduce is best used when you have a lot of discrete inputs, like webpages. We then apply `Map` to each input in parallel, then `Reduce` the results. It's typical to chain a bunch of these together, using outputs as inputs for other jobs.

## Word Count Example

![Alt text](image.png)