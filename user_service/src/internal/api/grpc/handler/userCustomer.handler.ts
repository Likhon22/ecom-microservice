import type { ServerUnaryCall, sendUnaryData } from '@grpc/grpc-js';
import type { UserCustomerService } from '../../../service/userCustomer.service.js';
import {
  CreateCustomerResponse,
  GetCustomersResponse,
  type CreateCustomerRequest,
  type GetCustomerByEmailRequest,
  type GetCustomersRequest,
} from '../../../../proto/gen/user_pb.js';
import { mapGrpcError } from '../utils/mapError.js';
import { toCreateCustomerResponse } from '../utils/responseFactory.js';

export class UserCustomerGrpcHandler {
  private service: UserCustomerService;
  constructor(service: UserCustomerService) {
    this.service = service;
  }
  async createCustomer(
    call: ServerUnaryCall<CreateCustomerRequest, unknown>,
    callback: sendUnaryData<CreateCustomerResponse>,
  ) {
    try {
      const result = await this.service.create({
        name: call.request.name,
        email: call.request.email,
        password: call.request.password,
        phone: call.request.phone,
        address: call.request.address,
        avatarUrl: call.request.avatarUrl,
      });
      const response = toCreateCustomerResponse(result);
      callback(null, response);
    } catch (error: unknown) {
      callback(mapGrpcError(error), null);
    }
  }
  async getCustomerByEmail(
    call: ServerUnaryCall<GetCustomerByEmailRequest, unknown>,
    callback: sendUnaryData<CreateCustomerResponse>,
  ) {
    try {
      const result = await this.service.getByEmail(call.request.email);
      const response = toCreateCustomerResponse(result);
      callback(null, response);
    } catch (err) {
      callback(mapGrpcError(err), null);
    }
  }
  async getCustomers(
    _call: ServerUnaryCall<GetCustomersRequest, unknown>,
    callback: sendUnaryData<GetCustomersResponse>,
  ) {
    try {
      const customers = await this.service.get();
      const response = new GetCustomersResponse({
        customers: customers.map(customer =>
          toCreateCustomerResponse(customer),
        ),
      });
      callback(null, response);
    } catch (err) {
      callback(mapGrpcError(err), null);
    }
  }
}
