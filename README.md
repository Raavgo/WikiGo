# WikiGo
Private Project to create a Wiki in Go

This repository contains:

- `frontend/` – a minimal [Astro](https://astro.build/) app.
- `backend/` – a Go server providing a JSON API.

## Backend

Generate an authentication token locally:

```bash
cd backend
go run token_gen.go > token.txt
export BLOG_TOKEN=$(cat token.txt)
```

Run the server:

```bash
go run .
```

### Endpoints
- `GET /api/hello` – sanity check.
- `POST /api/posts` – create a post from JSON (`title` and `content`). Requires `Authorization: Bearer $BLOG_TOKEN` header.
- `POST /api/posts/ipynb` – create a post from a Jupyter notebook. Send JSON `{ "title": "..", "notebook": "<ipynb as JSON>" }`.
- `GET /api/posts/list` – list all posts.

## Frontend

Run in development:

```bash
cd frontend
npm run dev
```
