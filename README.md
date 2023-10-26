# Abstract API
Abstract Golang API for MongoDB with Dynamic Collection and Advanced Searching.

## About
Originally conceived for personal use, this project was designed to facilitate internal testing without the need for building a full-fledged API. It empowers users to dynamically create collections and fields, adapting to the evolving needs of your data structures with ease.

## Dependecies
- Docker (with compose)

## Configuration
Create an `.env` file with application settings.

```
DB_HOST=db
DB_USER=root
DB_PASS=root
DB_NAME=data
DB_PORT=27017
API_PORT=8080
DBA_PORT=8081
```

## How to Run
Just run `docker-compose up -d` to create stack.

## Example of Use
This is an complete API flow in javascript.

```javascript
async function example() {
  var payload;

  // Create user
  payload = { name: "Administrator", email: "admin@admin.admin" };
  const user = await (await fetch(`http://localhost:8080/users`, { method: "POST", body: JSON.stringify(payload) })).json();

  // Create post
  payload = { title: "Post Title", content: "Post content text.", user_id: user.id };
  const post = await (await fetch(`http://localhost:8080/posts`, { method: "POST", body: JSON.stringify(payload) })).json();

  // Create comment
  payload = { content: "Post comment test.", post_id: post.id };
  const comment = await (await fetch(`http://localhost:8080/comments`, { method: "POST", body: JSON.stringify(payload) })).json();

  // Search posts
  const posts = await (
    await fetch(`http://localhost:8080/posts/search`, {
      method: "POST",
      body: JSON.stringify({
        paging: { page: 1, limit: 20 },
        filter: {
          published: true,
          // Filter posts with "post" on title
          title: { $regex: "post", $options: "i" },
        },
        join: {
          comments: "post_id",
          favorites: "post_id",
        },
        load: { users: "user:user_id" },
      }),
    })
  ).json();

  // Delete comment
  const resp = await (await fetch(`http://localhost:8080/comments/${comment.id}`, { method: "DELETE" })).json();

  // Update user
  payload = { name: "Administrator (edited)", email: "admin@admin.admin" };
  const edit = await (await fetch(`http://localhost:8080/users`, { method: "PUT", body: JSON.stringify(payload) })).json();
}
```

## Outros

- You can access an MongoDB GUI on port defined by `DBA_PORT`.
    - The credentials is `DB_USER` and `DB_PASS`.