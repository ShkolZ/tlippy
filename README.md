# Tlippy

Tlippy is a fast way to bulk download twitch clips by category or streamer. Tlippy is well-suited if
you don't want to go through each clip one by one to download it.

# Getting Started

To start using this cli app you need to setup you are going to need to register your own Twitch Application here https://dev.twitch.tv/.
Then you can get your Client-Id and Client-Secret and setup your .env like this:

```
CLIENT_ID="<your client id goes here>"
CLIENT_SECRET="<your client secret goes here>"
```

Then you're gonna need to build your go app

```
cd tlippy
go build . <custom-name>
```

with this you're set

# Tutorial

Basic usage is

```
<your-app-name> <path> <clip-amount>
tlippy ~/Videos/ 20
```

with those arguments you'are gonna download 20 clips most viewed clips of the months
