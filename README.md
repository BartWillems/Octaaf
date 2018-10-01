# Octaaf

A telegram bot written in Go

## Commands

* /aldi - Subscribe someone to the Aldi newsletter
* /avatar - Create an avatar
* /bodegem - A place that is real and exists
* /bol - Help someone out by subscribing them to the bol.com mailing lists!
* /care - Notify the participants that you don't give a hecky.
* /count - Get your room's current message count
* /doubt - When in doubt...
* /iasip - Get a random It's Always Sunny In Philadelphia quote
* /img - Search possible NSFW images
* /img_sfw - Search possible SFW images
* /issues - Show a list of open issues
* /kalirank - Show the kali rankings
* /m8ball - Let fate decide your future
* /mcaffee - Will McAffee eat his down dick?
* /next_launch - Show the next 5 rocket launches
* /quote - Get or store random kali quotes
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
1. a google api key *(optional)*
1. ```cp config/.env.dist config/.env```
1. ```cp config/database.yml.dist config/database.yml```
1. Enter the correct values in both config files
