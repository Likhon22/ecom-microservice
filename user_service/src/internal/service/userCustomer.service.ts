import type { UserCustomerRepo } from '../repo/userCustomer.repo.js';
import type {
  UserCustomerRequestDto,
  UserCustomerResponseDto,
} from '../domain/dtos/userCustmer.dto.js';
import type { TUser } from '../domain/user.domain.js';
import type { TCustomer } from '../domain/customer.domain.js';
import { hashPassword } from '../utils/hashPassword.js';
import ApiError from '../error/appError.js';
import mongoose from 'mongoose';

export interface UserCustomerServiceInterface {
  create(user: UserCustomerRequestDto): Promise<UserCustomerResponseDto>;
  get(): Promise<UserCustomerResponseDto[]>;
  getByEmail(email: string): Promise<UserCustomerResponseDto>;
}

export class UserCustomerService implements UserCustomerServiceInterface {
  private readonly repo: UserCustomerRepo;

  constructor(repo: UserCustomerRepo) {
    this.repo = repo;
  }

  async create(
    payload: UserCustomerRequestDto,
  ): Promise<UserCustomerResponseDto> {
    const session = await mongoose.startSession();
    session.startTransaction();

    try {
      const passwordHash = await hashPassword(payload.password);
      // --- Create User ---
      const user: Partial<TUser> = {
        email: payload.email,
        password: passwordHash,
        role: 'customer',
      };
      const createdUser = await this.repo.createUser(user, session);

      // --- Create Customer linked to User ---
      const customer: Partial<TCustomer> = {
        name: payload.name,
        email: payload.email,
        user: createdUser._id, // reference to user
        ...(payload.phone && { phone: payload.phone }),
        ...(payload.address && { address: payload.address }),
        ...(payload.avatarUrl && { avatarUrl: payload.avatarUrl }),
      };
      const createdCustomer = await this.repo.createCustomer(customer, session);

      // --- Commit transaction ---
      await session.commitTransaction();

      return {
        name: createdCustomer.name,
        email: createdCustomer.email,
        role: createdUser.role,
        status: createdUser.status,
        ...(createdCustomer.address && { address: createdCustomer.address }),
        ...(createdCustomer.avatarUrl && {
          avatarUrl: createdCustomer.avatarUrl,
        }),
        ...(createdCustomer.phone && { phone: createdCustomer.phone }),
      };
    } catch (error) {
      await session.abortTransaction();
      throw new ApiError(500, `User-Customer creation failed: ${error}`);
    } finally {
      session.endSession();
    }
  }
  async get(): Promise<UserCustomerResponseDto[]> {
    const customers = await this.repo.get();
    return customers.map(customer => ({
      name: customer.name,
      email: customer.email,
      role: customer.user.role,
      status: customer.user.status,
      ...(customer.phone && { phone: customer.phone }),
      ...(customer.avatarUrl && { avatarUrl: customer.avatarUrl }),
      ...(customer.address && { address: customer.address }),
    }));
  }
  async getByEmail(email: string): Promise<UserCustomerResponseDto> {
    const customer = await this.repo.getByEmail(email);
    if (!customer) {
      throw new ApiError(404, 'customer not found');
    }
    return {
      name: customer.name,
      email: customer.email,
      role: customer.user.role,
      status: customer.user.status,
      ...(customer.phone && { phone: customer.phone }),
      ...(customer.avatarUrl && { avatarUrl: customer.avatarUrl }),
      ...(customer.address && { address: customer.address }),
    };
  }
}
