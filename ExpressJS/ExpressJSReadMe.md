### Files included
I have only included the bare minimum. 
/app is the actual web-app
/API is the api

### Dependencies
node.js
postgresql

## How to run
#### Local development
1. Run 'npm install in root folder'
2. Populate ./utils/.env 
2. Navigate to /app for web-application-service or /API for the api-service
3. To run web-app: 'npm start', 
4. To run API: 'npm run start-api'



### How to query API
- **Get latest message ID**: localhost:5001/latest
- **Get all messages**: localhost:5001/msgs?latest=1
- **Get all messages for specific user**: localhost:5001/msgs/{username}?latest=1
- **Create message from user**: localhost:5001/msgs/{username}?latest=1
 - - body: {"content": "*message here*"}

- **Follow/Unfollow**: localhost:5001/fllw/{username}?latest=1

- - body: {"follow": {"username to follow"}}
- - body: {"unfollow": {"username to unfollow"}}

- **Register a user**: localhost:5001/register?latest=1
{
    "username": ...
    "email": ...
    "pwd": ...
}