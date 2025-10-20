/* eslint-disable @typescript-eslint/no-explicit-any */
import express, {
  type Application,
  type Request,
  type Response,
} from 'express';

import startServer from '../../internal/server.js';
import { connectDB } from '../../internal/infra/db/connection.js';

const app: Application = express();

app.use(express.json());
app.use(express.urlencoded());

app.get('/', (req: Request, res: Response) => {
  res.send('user service is running');
});

async function main() {
  try {
    const db = await connectDB();
    await startServer(app, db);
    console.log('Server is running with db');
  } catch (err: any) {
    console.log('Failed to start the server', err);
  }
}

main();
