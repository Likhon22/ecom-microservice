import type { ConnectRouter } from '@bufbuild/connect';
import { UserService } from '../../../../proto/gen/user_connect.js';
import type { UserCustomerGrpcHandler } from './userCustomer.grpc.js';

export function registerGrpcRoutes(
  router: ConnectRouter,
  grpcCustomerHandler: UserCustomerGrpcHandler,
) {
  // @ts-expect-error - ESM/CJS type mismatch in @bufbuild/protobuf
  router.service(UserService, {
    createCustomer:
      grpcCustomerHandler.createCustomer.bind(grpcCustomerHandler),
    getCustomerByEmail:
      grpcCustomerHandler.getCustomerByEmail.bind(grpcCustomerHandler),
    getCustomers: grpcCustomerHandler.getCustomers.bind(grpcCustomerHandler),
  });
}
