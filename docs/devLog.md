# Dev log

The purpose of this is to keep a 'blog' of updates by me and / or the rest of the project devs, as we learn networking & new areas of Go!

## 24/09/2023

I have tested out a simple proxy listener on the client, setting up the firefox SOCKS5 proxy and doing a simple net.listen in Go. When printing the object, it prints a memory location. Next I need to unpack whats actually there!

Although I haven't yet architected the project, I feel the starting point of the client code is good for a basic working proxy handler; before I implement too many features I do want the project architected at a high level. An infinite for loop runs which is constantly listening for stream events from the listener.