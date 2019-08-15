# n3-crdt
CRDT framework for n3 infrastructure

Send & Recive any json files, will be wrapped in version-vector crdt, sent to messaging server, and merged on receive with any other changes to the same objects received from other users.

For receive to work an instance of nats-streaming-server needs to be running on your machine.


