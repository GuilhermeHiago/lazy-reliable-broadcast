# Lazy Reliable Broadcast
Implementation of lazy reliable broadcast in Golang.

This implementations is made in modules, where each one has an *request* and *indication* chanel. There are 4 modules:
* *PP2PLink*: Perfect Point to Point link.
* *BEB*: Best Effort Broadcast.
* *PFD*: Perfect Failure Detector.
* *LRB*: Lazy Reliable Broadcast.
