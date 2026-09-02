# Todos

- [x] investigate how to properly handle password protected shares
- [ ] finish documenting `/shares` routes

## password protected shares

I'm unsure how i should structure the `GET /shares/:shareId` when it's password protected. Should the user send the cleartext password as a query param? Should instead there be a distinct endpoint that issues a JWT that allows the user to call `GET /shares/:shareId?jwt=<token>` and that verifies if the password is correct? Would both of these methods be insecure as the GET request caches the results? Investigate further.
