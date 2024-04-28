# **Auth app**
## - Access token type - **JWT**
## - Refresh token type - **JWT**
Refresh token stored in databse as **SHA512** hash (as jwt encryption algorythm) with GUID (to relate token to certain user)

Tokens are related to each other via creation time

---
## How to run this amazing repo:
1) read example of .env file ***(!!! be carefull, please, provide same internal and external port (it is required for stable application work) !!!)***</br>
1.1) also be careful and instead of localhost provide 0.0.0.0</br>
1.2) also be careful and dont forget to provide mongo url same as it is provided in docker-compose service name
2) write your own .env file (check out the docker-compose file and expose proper ports)
3) run: ```make``` or ```sudo make``` if ur docker daemon has no permissions
---
**API documentation is available on /swagger/index.html**