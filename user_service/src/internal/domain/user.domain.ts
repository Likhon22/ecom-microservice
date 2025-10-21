import type { Types } from 'mongoose';

export interface TUser {
  _id: Types.ObjectId;
  id: string;
  password: string;
  passwordChangedAt?: Date;
  email: string;
  role: 'admin' | 'customer' | 'superAdmin';
  status: 'in-progress' | 'blocked';
  isDeleted: boolean;
  createdAt: Date;
  updatedAt: Date;
}
export type TGender = 'male' | 'female' | 'other';
