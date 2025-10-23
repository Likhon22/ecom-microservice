import type {
  CreateCustomerRequest,
  CreateCustomerResponse,
  CustomerCredentialsResponse,
  DeleteCustomerRequest,
  DeleteCustomerResponse,
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
  DeleteCustomer: handleUnaryCall<
    DeleteCustomerRequest,
    DeleteCustomerResponse
  >;

  GetCustomerCredentials: handleUnaryCall<
    GetCustomerByEmailRequest,
    CustomerCredentialsResponse
  >;
};
