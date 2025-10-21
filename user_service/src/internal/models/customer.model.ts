import { Schema, model } from 'mongoose';
import type { TCustomer } from '../domain/customer.domain.js';

const customerSchema = new Schema<TCustomer>(
  {
    name: { type: String, required: true },
    phone: { type: String },
    user: { type: Schema.Types.ObjectId, ref: 'User', required: true },
    address: { type: String },
    avatarUrl: { type: String },
    email: { type: String, required: true, unique: true },
  },
  {
    timestamps: true,
  },
);

export const CustomerModel = model<TCustomer>('Customer', customerSchema);
