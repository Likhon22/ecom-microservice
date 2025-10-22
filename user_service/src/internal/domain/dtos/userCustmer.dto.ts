export interface UserCustomerRequestDto {
  name: string;
  phone?: string;
  address?: string;
  avatarUrl?: string;
  email: string;
  password: string;
}

export interface UserCustomerResponseDto {
  name: string;
  phone?: string;
  address?: string;
  avatarUrl?: string;
  email: string;
  role: 'admin' | 'customer' | 'superAdmin';
  status: 'in-progress' | 'blocked';
}
