# REAL-TIME-FORUM-TYPING-IN-PROGRESS

The goal of the project was to create a forum on which some of the data is updated in real time.

## Features

- View posts
  - Posts update real-time
- View post with it's comments
  - Comments update real-time
- Create new posts
- Create new comments
- View registered users
  - User status update real-time
  - Users are sorted by last activity, then username
- Chat with other users
  - Recieve notification on new message
  - Messages update real-time
  - Scroll to load more messages (throttled)
  - Display user typing status
- View users' profiles

## Implementation

- **SQLite3** DB
- **Go** backend (server, API, DB handling)
- **JS** front-end, app
- **HTML**, **CSS** styling

## Deployment

To run the server:

```bash
  go mod download
  go run .
```

To reset DB:

```bash
  go run . -db-reset
```

The app will run on: <http://localhost:8080/>

## Author

- [@Sagar Yadav](https://github.com/sagarishere)
- [@healingdrawing](https://healingdrawing.github.io/)
- [@acki02](https://01.gritlab.ax/git/acki02)

```

```
