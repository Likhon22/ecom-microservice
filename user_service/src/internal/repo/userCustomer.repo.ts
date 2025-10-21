import type { Model, ClientSession } from 'mongoose';
import type { TUser } from '../domain/user.domain.js';
import type { TCustomer } from '../domain/customer.domain.js';

export interface UserCustomerRepoInterface {
  createUser(user: Partial<TUser>, session?: ClientSession): Promise<TUser>;
  createCustomer(
    customer: Partial<TCustomer>,
    session?: ClientSession,
  ): Promise<TCustomer>;
}

export class UserCustomerRepo implements UserCustomerRepoInterface {
  private readonly userModel: Model<TUser>;
  private readonly customerModel: Model<TCustomer>;

  constructor(userModel: Model<TUser>, customerModel: Model<TCustomer>) {
    this.userModel = userModel;
    this.customerModel = customerModel;
  }

  async createUser(
    user: Partial<TUser>,
    session?: ClientSession,
  ): Promise<TUser> {
    const createdUser = await this.userModel.create([{ ...user }], { session });
    if (!createdUser[0]) throw new Error('User creation failed');
    return createdUser[0].toObject();
  }

  async createCustomer(
    customer: Partial<TCustomer>,
    session?: ClientSession,
  ): Promise<TCustomer> {
    const createdCustomer = await this.customerModel.create([{ ...customer }], {
      session,
    });
    if (!createdCustomer[0]) throw new Error('Customer creation failed');
    return createdCustomer[0].toObject();
  }
}
