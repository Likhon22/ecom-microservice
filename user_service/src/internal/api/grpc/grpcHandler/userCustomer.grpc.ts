/* eslint-disable @typescript-eslint/no-unused-vars */
import type {
  CreateCustomerRequest,
  CreateCustomerResponse,
  GetCustomerByEmailRequest,
  GetCustomersRequest,
  GetCustomersResponse,
} from '../../../../proto/gen/user_pb.js';
import type { UserCustomerRequestDto } from '../../../domain/dtos/userCustmer.dto.js';
import type { UserCustomerService } from '../../../service/userCustomer.service.js';
import { handleGrpcError } from '../../middleware/grpcErrorHanlder.js';

export class UserCustomerGrpcHandler {
  private readonly service: UserCustomerService;
  constructor(service: UserCustomerService) {
    this.service = service;
  }
  async createCustomer(
    req: CreateCustomerRequest,
  ): Promise<CreateCustomerResponse> {
    try {
      console.log('i am there');

      const dto: UserCustomerRequestDto = {
        name: req.name,
        email: req.email,
        password: req.password,
        ...(req.phone && { phone: req.phone }),
        ...(req.address && { address: req.address }),
        ...(req.avatarUrl && { avatarUrl: req.avatarUrl }),
      };
      const result = await this.service.create(dto);
      return result as unknown as CreateCustomerResponse;
    } catch (error) {
      throw handleGrpcError(error);
    }
  }
  async getCustomerByEmail(
    req: GetCustomerByEmailRequest,
  ): Promise<CreateCustomerResponse> {
    try {
      const result = await this.service.getByEmail(req.email);
      return result as unknown as CreateCustomerResponse;
    } catch (error) {
      throw handleGrpcError(error);
    }
  }

  async getCustomers(req: GetCustomersRequest): Promise<GetCustomersResponse> {
    try {
      const list = await this.service.get();
      return { customers: list } as unknown as GetCustomersResponse;
    } catch (error) {
      throw handleGrpcError(error);
    }
  }
}
