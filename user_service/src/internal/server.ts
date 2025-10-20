import type { Application } from 'express';
import type { Server } from 'http';
import config from './config/index.js';
import type { Connection } from 'mongoose';

let httpServer: Server;

async function startServer(app: Application, db: Connection): Promise<Server> {
  try {
    httpServer = app.listen(config.port, () => {
      console.log(`User service listening on port ${config.port}`);
    });
    return httpServer;
  } catch (err) {
    console.error(err);
    throw err;
  }
}

// Handle unhandled rejections
process.on('unhandledRejection', (reason, err) => {
  console.error('Unhandled Rejection caught. Shutting down...', reason, err);
  if (httpServer) {
    httpServer.close(() => process.exit(1));
  } else {
    process.exit(1);
  }
});

// Handle uncaught exceptions
process.on('uncaughtException', err => {
  console.error('Unhandled Exception caught. Shutting down...', err);
  process.exit(1);
});

export default startServer;
