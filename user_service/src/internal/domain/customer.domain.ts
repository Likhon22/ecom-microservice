import type { TUser } from './user.domain.js';

export interface TCustomer extends TUser {
  firstName: string;
  lastName: string;
  phone?: string;
  address?: string;
  avatarUrl?: string;
}
