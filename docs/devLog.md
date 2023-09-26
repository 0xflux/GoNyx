# Dev log

The purpose of this is to keep a 'blog' of updates by me and / or the rest of the project devs, as we learn networking & new areas of Go!

## 24/09/2023

I have tested out a simple proxy listener on the client, setting up the firefox SOCKS5 proxy and doing a simple net.listen in Go. When printing the object, it prints a memory location. Next I need to unpack whats actually there!

Although I haven't yet architected the project, I feel the starting point of the client code is good for a basic working proxy handler; before I implement too many features I do want the project architected at a high level. An infinite for loop runs which is constantly listening for stream events from the listener.

## 25/09/2023

The proxy server now properly handles a SOCKS5 standard request, thanks to some well documented resources on Google as to SOCKS5 handshakes. At the moment we are simply intercepting the domain and port; there is no availability for URI inspection, so that will likely be the next area of focus. Some screenshot examples of whats going on under the hood.

I'm not sure what DNS is doing though, it is set to route via the SOCKS5 proxy - a request to example.nyx tries to make a Google search (that's what the google:443 request is in the below screenshot). I'm not sure where this routing logic is, whether it is a DNS request not finding *.nyx, or whether it's some functionality baked into Firefox.

![image](https://github.com/the-wandering-photon/GoNyx/assets/49762827/7e457416-2a56-4f9d-9aa0-f0a70c466943)

## 26/09/2023

Spent quite some time today looking into the next part of the connection stream, which should be the HTTP data the proxy processes, after validating the SOCKS5 handshake. I currently expect to be able to read into a buffer the http connection data; however it appears to be empty (EOF). More granular debugging is required, there could be data there but I am accessing it wrong, or the browser closed the stream before the HTTP data could be processed (unlikely), or there is some issue with the connection state. Next steps will be to debug this. As far as I can see from SOCKS5 documentation, I do not need to send a response back to the client until the connection was successfully forwarded to the destination, so I don't think the HTTP data would follow if I send a response to the client - however I am remaining open minded to this.
