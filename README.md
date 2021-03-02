# requestpq
A minimal implementation of priority queue for web requests. 

Concurrency may not be as nice as it seems when we are serving deep models at the backend. These models usually have large FLOPs and consumes high utilization of CPU/GPU as well as memory usage. To exploit the hardware capability and avoid OOM, a better way is to establish a queue, for which CPU/GPU workers are the consumers and request handlers are the producers.

In some scenarios, requests do not have the same weights. Considering the resource constraints, we hope to serve tasks with higher priority first to reduce their latency, thus here is the minimal solution for it! 