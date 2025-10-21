import bcrypt from 'bcrypt';
import config from '../config/index.js';

export async function hashPassword(password: string): Promise<string> {
  return await bcrypt.hash(password, Number(config.salt_round));
}
