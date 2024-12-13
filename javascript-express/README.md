# Express JS Application

## Files included

I have only included the bare minimum.
/app is the actual web-app
/API is the api

### Dependencies

node.js
postgresql

## How to run

### Local development

1. Run 'npm install in codebase/ExpressJS'
2. Make sure to populate .env in root folder
3. To run web-app: 'npm start', at codebase/ExpressJS
4. To run API: 'npm run start-api' at codebase/ExpressJS

### How to query API

- **Get latest message ID**: GET - localhost:5001/latest
- **Get all messages**: GET - localhost:5001/msgs?latest={provide a number}
- **Get all messages for specific user**: GET - localhost:5001/msgs/{username}?latest={provide a number}
- **Get all followers of a user**: GET - localhost:5001/fllws/{username}?latest={provide a number}
- **Create message from user**: POST - localhost:5001/msgs/{username}?latest={provide a number}
  -- body: {"content": "_your message here_"}

- **Follow/Unfollow**: localhost:5001/fllw/{username}?latest={provide a number}

    -- body: {"follow": "username to follow"}

    -- body: {"unfollow": {"username to unfollow"}}

- **Register a user**: localhost:5001/register?latest={provide a number}
  -- body: {
  "username": "...",
  "email": "...",
  "pwd": "..."
  }
