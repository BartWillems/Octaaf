# Octaaf

A telegram bot written in Go

## Commands

* /all - Send a message to all active members in a group
* /bodegem - A place that is real and exists
* /care - Notify the participants that you don't give a hecky.
* /changelog - View the Octaaf changelog
* /count - Get your room's current message count
* /doubt - When in doubt...
* /iasip - Get a random It's Always Sunny In Philadelphia quote
* /img - Search possible NSFW images
* /img_sfw - Search possible SFW images
* /kaleaderboard - Show the kaleets or blazers
* /kalirank - Show the kali rankings
* /m8ball - Let fate decide your future
* /more - MORE IMAGES
* /next_launch - Show the next 5 rocket launches
* /pollentiek - Shows your political orientation, by doing machine learning, AI and blockchain in the cloud with microservices.
* /quote - Get or store random kali quotes
* /remind_me - Remind me in a given time
* /roll - Praise kek
* /search - Search stuff on DuckDuckGo with safe search on
* /search_nsfw - Search dirty stuff on DuckDuckGo
* /stallman - I'd just like to interject for a moment. What you’re referring to as Linux, is in fact, GNU/Linux, or as I’ve recently taken to calling it, GNU plus Linux.
* /weather - Get the weather of a city
* /what - Explains what something is
* /where - Find places on earth
* /xkcd - Get a random XKCD comic

## Developing

### Requirements

1. a telegram bot account
    * you can use telegram for this
1. install postgresql
    * ```pacman -S postgresql``` (or any other package manager)
    * ```sudo -u postgres -i```
    * ```initdb --locale en_US.UTF-8 -E UTF8 -D '/var/lib/postgres/data'```
    * [Create your first DB user](https://wiki.archlinux.org/index.php/PostgreSQL#Create_your_first_database.2Fuser)
1. Install and run Redis
1. [Soda](https://gobuffalo.io/en/docs/db/toolbox)
1. a google api key *(optional)*
1. ```cp config/settings.toml.dist config/settings.toml```
1. Enter the correct values in the settings file
1. _(Optional)_ <https://github.com/tools/godep>

#### But I don't know how to computer 😨😨😨

1. ```cp config/settings.toml.dist config/settings.toml```
1. Enter the correct values in config/settings.toml (just enter your telegram api key)
1. ```docker-compose up```
