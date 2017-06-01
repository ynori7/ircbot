IRC Bot
======

A simple Golang IRC bot. It handles all the communication with the IRC server, parsing all output from the server and sending commands back to the server.

Features
-------
Currently the bot supports the following features:

- Greets users as they join a channel which the bot is in
- Responds to greetings directed towards the bot
- Automatically rejoins channels when kicked from them
- Looks up Github repos from the Github API and shows info about it
- Automatically gives voice to users who join a channel that the bot moderates, and maintains a list of muted users to not give voice to (controlled using the mute/unmute commands)

Usage
----
To run the IRC bot, checkout the repo, install the dependencies, build the code, and then execute, passing in the path to your config file. For example:

`%GOPATH%/bin/ircbot.exe "path/to/config.yml"`


Configuration
-------------
The config file should have the following:

- connection_string: The server and port, for exmaple: "irc.psych0tik.net:6697"
- nick: The name of the bot
- channels: The list of channels which the bot should join
- moderated_channels: A list of channels which are "moderated" (having the +m mode). The bot (assuming it has the proper privileges) will give "voice" (+v) to users when they join these channels.
- password: The password used to identify with NickServ (ignored if empty)
- admins: A list of users who are considered to be administrators for this bot. Certain commands may only be available for admin users.
- use_ssl: Either true or false if the server requires SSL
- greetings: A list of greetings the bot should respond to (e.g. "hi", "hello", "sup", etc.)
