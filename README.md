# go_chat
Simple GO chat service using websocket and grpc

There are 3 components: a `message store`, a `webserver` and a simple `frontend`.
The `webserver` listens for incoming http connections from the `frontend` and tries to update it to a web socket connection. When a new connection comes in, it will be registered and messages from different clients will be broadcasted between them. The `message store` component is optional, but if it does operate then on a new clients connection it will receive all the previous messages made by other clients. Currently ports are hardcoded in all 3 components: `webserver` listens on `8080` while the `message store` on `8081`.

To start the backend components, execute the following commands in two different terminal:
 - `go run .\cmd\grpc_server\main.go` to start the `message store`,
 - `go run .\cmd\ws_server\main.go` to start the `webserver`,
and open the `frontend\home.html`  on multiple pages in your web sockets enabled browser.
