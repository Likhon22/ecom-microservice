/* eslint-disable @typescript-eslint/no-explicit-any */
import express, {
  type Application,
  type Request,
  type Response,
} from 'express';
import { connectNodeAdapter } from '@bufbuild/connect-node';
import startServer from '../../internal/server.js';
import { connectDB } from '../../internal/infra/db/connection.js';
import { UserCustomerRepo } from '../../internal/repo/userCustomer.repo.js';
import { UserModel } from '../../internal/models/user.model.js';

import { UserCustomerHandler } from '../../internal/api/handlers/userCustomer/userCustomer.handler.js';
import { UserCustomerService } from '../../internal/service/userCustomer.service.js';
import { registerRoutes } from '../../internal/api/routes/index.js';
import { Middleware } from '../../internal/api/middleware/index.js';
import { CustomerModel } from '../../internal/models/customer.model.js';
import { UserCustomerGrpcHandler } from '../../internal/api/grpc/grpcHandler/userCustomer.grpc.js';
import { registerGrpcRoutes } from '../../internal/api/grpc/grpcHandler/grpcRoutes.js';

const app: Application = express();

app.get('/', (req: Request, res: Response) => {
  res.send('user service is running');
});

async function main() {
  try {
    //connect to database
    await connectDB();

    // middleware
    const mw = new Middleware();

    //repo
    const userCustomerRepo = new UserCustomerRepo(UserModel, CustomerModel);

    //service
    const userCustomerService = new UserCustomerService(userCustomerRepo);

    //handler
    const userCustomerHandler = new UserCustomerHandler(userCustomerService);

    //grpc handler

    const userCustomerGrpcHandler = new UserCustomerGrpcHandler(
      userCustomerService,
    );

    //grpc routes
    app.use(
      '/grpc',
      connectNodeAdapter({
        routes: router => registerGrpcRoutes(router, userCustomerGrpcHandler),
      }),
    );
    app.use(express.json({ type: 'application/json' }));
    app.use(express.urlencoded({ extended: true }));
    // rest routes
    registerRoutes(app, userCustomerHandler, mw);
    // global error handler
    app.use(mw.globalErrorHandler);

    // not found routes
    app.use(mw.noRoutesFound);

    await startServer(app);
    console.log('Server is running with db mongodb');
  } catch (err: any) {
    console.log('Failed to start the server', err);
  }
}

main();
