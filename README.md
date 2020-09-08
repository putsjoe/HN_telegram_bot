# HN_telegram_bot
A Telegram bot written in Go to view and save the latest HackerNews posts (Golang, SQLite).

I didnt want to rely on my browser for new hacker news posts anymore, instead this provides a way of saving posts to a database that I want to read later.

Upon running the service, it will scan and grab the top 500 posts from HackerNews and save them. When the user sends the `/latest` command to the bot, five unread posts will be shown and then marked as read in the database. Every hour the service gets the top 500 posts and saves any with a score over 10 that havent been saved yet.

It may be a good idea to access the database and mark all as read, otherwise the first time its used you will have 500 posts to go through.


## Setup

First setup a new bot on Telegram, instructions for this are [here](https://core.telegram.org/bots#creating-a-new-bot) which will give you a bot API key to be used in the blow instructions. Run the following from the project root:

```
go build main.go
sqlite3 data.sqlite < hackernews/schema.sql

./main -key={{ your bot API key }}
```
You should then be able to run the keyboard commands to start using the project.

### Restrict to user
An optional `-user` flag can be passed to restrict the bot to only send messages to the given user id. To find your userID, use the @userinfobot. 

### Ping
Sending the `/ping` command to the bot returns the number of posts in the database and the number that havent been read yet.

## Limitations

- The first run has 500 posts marked as unread. Needs a command to mark all as read or a flag to pass to the service.
- Any time unread posts are returned, they are marked as read, therefore if there is more than one user they wont see all new. Needs support for multi user read.
- Upon saving or deleting a post, the command is entered and so scrolls to the bottom. Ideally a system to save and delete without having to enter a command but this may not be possible with Telegram.
- Golang code improvements - Tests, naming, packaging

### Further Notes
- An API key is exposed in old commits, however this bot has now been removed and so the API key will not work.

