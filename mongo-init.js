// mongo-init.js
db.createUser({
    user: process.env.MONGO_USER || 'root',
    pwd: process.env.MONGO_PASSWORD || 'example',
    roles: [
      { role: 'readWrite', db: process.env.DB_NAME || 'go_3_mini' }
    ]
  });