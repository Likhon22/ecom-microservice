import type { ServerUnaryCall, sendUnaryData } from '@grpc/grpc-js';
import type { UserCustomerService } from '../../../service/userCustomer.service.js';
import {
  CreateCustomerResponse,
  CustomerCredentialsResponse,
  DeleteCustomerRequest,
  DeleteCustomerResponse,
  GetCustomersResponse,
  type CreateCustomerRequest,
  type GetCustomerByEmailRequest,
  type GetCustomersRequest,
} from '../../../../proto/gen/user_pb.js';
import { mapGrpcError } from '../utils/mapError.js';
import { toCreateCustomerResponse } from '../utils/responseFactory.js';
import { TCustomerSchema } from '../../handlers/userCustomer/userCustomer.validation.js';

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
      TCustomerSchema.shape.body.parse(call.request);
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
  async deleteCUstomer(
    call: ServerUnaryCall<DeleteCustomerRequest, unknown>,
    callback: sendUnaryData<DeleteCustomerResponse>,
  ) {
    try {
      const response = await this.service.deleteUser(call.request.email);
      const msg = new DeleteCustomerResponse({
        msg: response,
      });
      callback(null, msg);
    } catch (error) {
      callback(mapGrpcError(error), null);
    }
  }
  async getCustomerCredentials(
    call: ServerUnaryCall<GetCustomerByEmailRequest, unknown>,
    callback: sendUnaryData<CustomerCredentialsResponse>,
  ) {
    try {
      const credentials = await this.service.getCustomerCredentials(
        call.request.email,
      );

      return callback(null, credentials);
    } catch (err) {
      callback(mapGrpcError(err), null);
    }
  }
}
