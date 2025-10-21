import mongoose from 'mongoose';
import type { UserCustomerRepo } from '../repo/userCustomer.repo.js';
import type { UserCustomerDto } from '../domain/dtos/userCustmer.dto.js';
import type { TUser } from '../domain/user.domain.js';
import type { TCustomer } from '../domain/customer.domain.js';
import { hashPassword } from '../utils/hashPassword.js';

export interface UserCustomerServiceInterface {
  create(user: UserCustomerDto): Promise<UserCustomerDto>;
}

export class UserCustomerService implements UserCustomerServiceInterface {
  private readonly repo: UserCustomerRepo;

  constructor(repo: UserCustomerRepo) {
    this.repo = repo;
  }

  async create(payload: UserCustomerDto): Promise<UserCustomerDto> {
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
      if (!createdUser) {
        throw new Error('User creation failed. Cannot create customer.');
      }
      if (!createdUser._id) {
        throw new Error('User creation failed. no id available.');
      }

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

      // --- Return safe DTO ---
      return {
        name: createdCustomer.name,
        email: createdCustomer.email,
        password: '',
      };
    } catch (error) {
      await session.abortTransaction();
      throw new Error(`User-Customer creation failed: ${error}`);
    } finally {
      session.endSession();
    }
  }
}
