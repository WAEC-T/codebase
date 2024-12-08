# DB-operations change

1. MyTimeline should be two not three ops. _GetMyMessages_ should add a sub-query for followers
2. FollowUser shoudl be two and not three. _FollowUser_ does not need to fetch the username of the logged in user. Instead, simply take the userid directly from the session.
3. ^same for UnFollow
4. RegisterUser do not need to get the GetUser logged in
5. AddMessage change from 2 to 1 ops. No need to get username in _GetUser_