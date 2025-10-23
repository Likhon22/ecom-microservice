import type { Types } from 'mongoose';
import type { TUser } from './user.domain.js';

export interface TCustomer {
  name: string;
  user: Types.ObjectId;
  phone?: string;
  address?: string;
  avatarUrl?: string;

  email: string;
}
export interface TCustomerPopulated extends Omit<TCustomer, 'user'> {
  user: Pick<TUser, 'role' | 'status' | 'isDeleted'>;
}
