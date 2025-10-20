import type { TGender, TUser } from './user.domain.js';

export interface TAdmin extends TUser {
  firstName: string;
  lastName: string;
  gender: TGender;
}
