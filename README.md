# Pond

Development version. Use with care.

## Dependencies

redis-server
go

## Usage

Expects a `REDIS_URL` and accpets a `PORT` to define the port to be binded

## Commands

```bash
# Starts a server and adds a list to ponds to send the messages
pond server -pond [friend pond] -pond [other friend]

# Fetches all the available messages from a pond
pond fetch -pond [the pond to be fetched and attempted to decrypt]

# Sends an encrypted and anonymous message to a pond
pond send -pond [where to send the message] email "Message"
```
