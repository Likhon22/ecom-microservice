export interface TUser {
  id: string;
  password: string;
  passwordChangedAt?: Date;
  email: string;
  role: 'admin' | '' | 'superAdmin';
  status: 'in-progress' | 'blocked';
  isDeleted: boolean;
  createdAt: Date;
  updatedAt: Date;
}
export type TGender = 'male' | 'female' | 'other';
