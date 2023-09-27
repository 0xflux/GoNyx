# High level topology

At the high level, a browser proxy is in place which connects the browser to the NyxNet, and traffic is routed through this to the RVP, where data is exchanged from the inbound connection from the intended .nyx domain.

The Nyx Client adn Nyx Server Client will poll the Accord to establish their routes, and the RVP is selected by the client at random. Whilst the client is running, it will maintain a route (not as a live conneciton).

The client & server will then handle the layered encryption from each node within the NyxNet, and then forward traffic through the NyxNet where it eventually meets at the RVP. The client connection is then received at the RVP, at 
which point the RVP orchestrates the request from the .nyx server, where that data is sent to the RVP. The client is waiting for the connection, which will be managed through internal ports, once data is exchanged it will then be routed back through NyxNet,
 again applying encryption, and gets shown in the browser.

![image](https://github.com/the-wandering-photon/GoNyx/assets/49762827/4baa0e7c-92f6-4b5f-b75f-416e9c9b9056)
