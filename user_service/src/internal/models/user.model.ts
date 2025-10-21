import { Schema, model } from 'mongoose';
import type { TUser } from '../domain/user.domain.js';

const userSchema = new Schema<TUser>(
  {
    password: { type: String, required: true },
    passwordChangedAt: { type: Date },
    email: { type: String, required: true, unique: true },
    role: {
      type: String,
      enum: ['admin', 'customer', 'superAdmin'],
      default: 'customer',
      required: true,
    },
    status: {
      type: String,
      enum: ['in-progress', 'blocked'],
      default: 'in-progress',
    },
    isDeleted: { type: Boolean, default: false },
  },
  {
    timestamps: true, // Automatically adds createdAt and updatedAt
  },
);

export const UserModel = model<TUser>('User', userSchema);
