# proLogChat
A host-client chat program that allows users to communicate with each other in a chat room-like format. First user to run program becomes the host. Each user after that becomes a client. They connect and send messages to host user. Host user sends all received messages to each connected user

For Host:
  The first user to open program automatically becomes the host
  Host listens for client connections, accepts messages, and then repeats all messages to every connected client
  
For Client:
  Every user to open program after the first becomes a client.
  Client will be immediately prompted for a username
  This username will accompany any message the user enters afterwards
  Client can then enter in the console any message followed by Enter and it will be sent to the Host
  The host will then distribute the message to all connected clients
