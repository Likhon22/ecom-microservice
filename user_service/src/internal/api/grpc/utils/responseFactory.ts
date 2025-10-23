import { CreateCustomerResponse } from '../../../../proto/gen/user_pb.js';
import type { UserCustomerResponseDto } from '../../../domain/dtos/userCustmer.dto.js';

export function toCreateCustomerResponse(
  dto: UserCustomerResponseDto,
): CreateCustomerResponse {
  return new CreateCustomerResponse({
    name: dto.name,
    email: dto.email,
    role: dto.role,
    status: dto.status,
    phone: dto.phone ?? '',
    address: dto.address ?? '',
    avatarUrl: dto.avatarUrl ?? '',
    isDeleted: dto.isDeleted,
  });
}
