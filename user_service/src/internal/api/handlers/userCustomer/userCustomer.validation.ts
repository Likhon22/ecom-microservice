import { z } from 'zod';

export const TCustomerSchema = z.object({
  body: z.object({
    name: z.string().min(1, 'Name is required'),
    phone: z.string().optional(),
    address: z.string().optional(),
    avatarUrl: z.string().url().optional(),
    email: z.string().email(),
    password: z.string().min(6, 'Password must be at least 6 characters'),
  }),
});

export type TUserZod = z.infer<typeof TCustomerSchema>;
