## What's this?
A set of tools to simulate a cloud architecture. The building blocks are

* instances, devices with a task to run, disk, and memory;
* groups, instances performing the same task and load balancing among 
  them;
* queues, like SQS; and
* buckets, blob storage.

Each building block have constraints. Instances have limited disk and
memory and only one task can run at a time. Buckets and queues have
worse response times and limited throughput to a single instance but
provide infinite storage and are not slowed down when many instances
are communicating with them simultaneously.

Let's say you use one instance to store results. If you have a group
of instances trying to access this instance their requests will be
processed one at a time. If you use a bucket, on the other hand, all
instances will be able to issue requests at the same time.

# Why?
With these blocks it is possible to build workflows with multiple
load-balanced groups of instances and queues that would otherwise take
a long time to implement, deploy, and make changes to.

It's also possible to build storage solutions that illustrate concepts
like sharding and replication, both on disk or in-memory storage
similar to Memcached.

## Future improvements
* Examples (see the tests for now).
* Configurable latencies.
* Add failures: disk, instances, networks.
* Visualize latencies: "CPU" and IO wait on instances, network 
  throughput, and so on.
* Add a relational database building block.
