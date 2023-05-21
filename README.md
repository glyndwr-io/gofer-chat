# GoferChat

GoferChat is a websocket based chat server written in Go with an included and optional Svelte Front-End.

## Running GoferChat

If you wish to use the built in GoferChat frontend written in Sevlte, you need to first complie the Svelte code. To do so, enter the following commands:

```
$ cd client
$ npm install
$ npm run build
$ cd ..
```

If you wish to build your own Front-End for GoferChat, place the files to be served in the `./client/build/` directory. Make sure the main application is named `index.html` and the login page is `login.html`. For any other static files, GoferChat will serve them statically for you in the `./client/build/` directory.

Once you've initalized the Front-End, you can build and run the backend by running the following commands:

```
$ go build .
$ ./gofer-chat
```
