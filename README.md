# Toussaint

[![Build Status](https://travis-ci.org/wuzzapcom/toussaint.svg?branch=master)](https://travis-ci.org/wuzzapcom/toussaint)

Toussaint is a project for tracking sales in Playstation Store.  
It provides an universal backend with single frontend telegram bot.  
But new frontends will appear in future.

## Deploy

### Preparations

We should prepare file with telegram token. Ensure that you in toussaint directory.
Create file with secret:

```
vim docker/telegram/telegram.token
```

Build and run:

```
docker-compose -f docker/telegram/telegram-compose.yml build
docker-compose -f docker/telegram/telegram-compose.yml up
```