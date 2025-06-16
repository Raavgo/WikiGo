# WikiGo
Private Project to create a Wiki in Go

This repository contains:

- `frontend/` – a minimal [Astro](https://astro.build/) app.
- `backend/` – a Go server providing a JSON API.

## Backend

Generate an authentication token locally:

```bash
cd backend
go run ./cmd/token > token.txt
export BLOG_TOKEN=$(cat token.txt)
```

Set the token used to unlock premium content (optional):

```bash
export PAYWALL_TOKEN=mysecret
```

Use this value when requesting premium content via `X-Paywall-Token` header.

Run the server:

```bash
go run .
```

### Endpoints
- `GET /api/hello` – sanity check.
- `POST /api/posts` – create a text post. Requires `Authorization: Bearer $BLOG_TOKEN`.
- `POST /api/posts/ipynb` – create a text post from a Jupyter notebook.
- `POST /api/posts/video` – upload a video file (`file`, `title`, `premium`).
- `POST /api/posts/audio` – upload an audio file (`file`, `title`, `premium`).
- `GET /api/posts/list?premium=false` – list free posts.
- `GET /api/posts/list?premium=true` – list premium posts (needs `X-Paywall-Token`).
- `GET /api/posts/media/{id}` – download media for a post.

## Frontend

Run in development:

```bash
cd frontend
npm run dev
```

The frontend provides a simple layout with sections for free and premium posts. To view premium posts, enter the paywall token when prompted.
