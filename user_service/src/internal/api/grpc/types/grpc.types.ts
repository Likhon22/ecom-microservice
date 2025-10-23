import type {
  CreateCustomerRequest,
  CreateCustomerResponse,
  GetCustomerByEmailRequest,
  GetCustomersRequest,
  GetCustomersResponse,
} from '../../../../proto/gen/user_pb.js';
import type { handleUnaryCall } from '@grpc/grpc-js';
export type UserServiceHandlers = {
  CreateCustomer: handleUnaryCall<
    CreateCustomerRequest,
    CreateCustomerResponse
  >;
  GetCustomerByEmail: handleUnaryCall<
    GetCustomerByEmailRequest,
    CreateCustomerResponse
  >;
  GetCustomers: handleUnaryCall<GetCustomersRequest, GetCustomersResponse>;
};
