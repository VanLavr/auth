# **Auth app**
## - Access token type - **JWT**
## - Refresh token type - **JWT**
Refresh token stored in databse as **SHA512** hash (as jwt encryption algorythm) with GUID (to relate token to certain user)

Tokens are related to each other via creation time

---
# How to run this amazing repo:
1) read example of .env file
2) write your own .env file
3) run: ```make```