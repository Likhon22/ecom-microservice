import type { Application } from 'express';
import type { Server } from 'http';
import config from './config/index.js';
import type { UserCustomerService } from './service/userCustomer.service.js';
import * as grpc from '@grpc/grpc-js';
import { UserCustomerGrpcHandler } from './api/grpc/handler/userCustomer.handler.js';
import { fileURLToPath } from 'url';
import path from 'path';
import { loadSync } from '@grpc/proto-loader';
import type { UserServiceHandlers } from './api/grpc/types/grpc.types.js';

let httpServer: Server;

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const PROTO_PATH = path.resolve(__dirname, '../proto/user.proto');

const packageDefinition = loadSync(PROTO_PATH, {
  keepCase: false,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
  includeDirs: [
    path.resolve(__dirname, '../proto'),
    path.resolve(__dirname, '../proto/validate'),
  ],
});

const proto = grpc.loadPackageDefinition(packageDefinition) as {
  user_service: {
    UserService: {
      service: unknown;
    };
  };
};
export async function startGrpcServer(service: UserCustomerService) {
  const server = new grpc.Server();
  const handler = new UserCustomerGrpcHandler(service);
  const userServiceDefinition = proto.user_service.UserService
    .service as grpc.ServiceDefinition<UserServiceHandlers>;
  server.addService(userServiceDefinition, {
    CreateCustomer: handler.createCustomer.bind(handler),
    GetCustomerByEmail: handler.getCustomerByEmail.bind(handler),
    GetCustomers: handler.getCustomers.bind(handler),
    DeleteCustomer: handler.deleteCUstomer.bind(handler),
    GetCustomerCredentials: handler.getCustomerCredentials.bind(handler),
  });
  const address = `0.0.0.0:${config.grpc_port ?? '5001'}`;
  await new Promise<void>((resolve, reject) => {
    server.bindAsync(address, grpc.ServerCredentials.createInsecure(), err => {
      if (err) {
        reject(err);
        return;
      }
      resolve();
    });
  });
  server.start();
  console.log(`grpc server running on port ${address}`);

  return server;
}

async function startServer(app: Application): Promise<Server> {
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
