import type { Model, ClientSession } from 'mongoose';
import type { TUser } from '../domain/user.domain.js';
import type {
  TCustomer,
  TCustomerPopulated,
} from '../domain/customer.domain.js';
import ApiError from '../error/appError.js';

export interface UserCustomerRepoInterface {
  createUser(user: Partial<TUser>, session?: ClientSession): Promise<TUser>;
  createCustomer(
    customer: Partial<TCustomer>,
    session?: ClientSession,
  ): Promise<TCustomer>;
  get(): Promise<TCustomerPopulated[]>;
  getByEmail(id: string): Promise<TCustomerPopulated | null>;
  deleteUser(email: string): Promise<TUser | null>;
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
    if (!createdUser[0]) throw new ApiError(500, 'User creation failed');
    return createdUser[0].toObject();
  }

  async createCustomer(
    customer: Partial<TCustomer>,
    session?: ClientSession,
  ): Promise<TCustomer> {
    const createdCustomer = await this.customerModel.create([{ ...customer }], {
      session,
    });
    if (!createdCustomer[0])
      throw new ApiError(500, 'Customer creation failed');
    return createdCustomer[0].toObject();
  }
  async get(): Promise<TCustomerPopulated[]> {
    const customers = await this.customerModel
      .find()
      .populate('user', 'role status')
      .lean<TCustomerPopulated[]>();
    return customers;
  }
  async getByEmail(email: string): Promise<TCustomerPopulated | null> {
    const customer = await this.customerModel
      .findOne({ email: email })
      .populate('user', 'role status isDeleted')
      .lean<TCustomerPopulated>();
    return customer;
  }
  async deleteUser(email: string): Promise<TUser | null> {
    const user = await this.userModel
      .findOneAndUpdate({ email }, { isDeleted: true }, { new: true })
      .lean<TUser>();
    return user;
  }
}
