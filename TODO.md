# 📋 TODO

- [x] Add a `TODO.md` file to the project.
- [ ] Re-implement MongoDB.
  - [ ] Add option to delete chat history at end of conversation.
- [x] Insert single prompt and response endpoint without db
- [x] Implement caching
- [x] Implement HTTP server with endpoints for:
  - [x] `/chat` - POST
  - [ ] `/chat/:id` - GET
  - [x] `/prompt` - POST
- [ ] Shell should use the HTTP server; OR have a local implementation and a server implementation. See what's better
- [ ] Success message for connection to DB should check DB first