import type { Request, Response } from 'express';

import type { UserCustomerService } from '../../../service/userCustomer.service.js';
import type { UserCustomerRequestDto } from '../../../domain/dtos/userCustmer.dto.js';
import sendResponse from '../../../utils/sendResponse.js';

export class UserCustomerHandler {
  private readonly service: UserCustomerService;
  constructor(service: UserCustomerService) {
    this.service = service;
  }

  async create(req: Request, res: Response) {
    const payload = req.body as UserCustomerRequestDto;

    const user = await this.service.create(payload);
    sendResponse(res, {
      message: 'user created successfully',
      data: user,
      statusCode: 200,
    });
  }
  async get(req: Request, res: Response) {
    const customers = await this.service.get();
    sendResponse(res, {
      message: 'customers fetched successfully',
      data: customers,
      statusCode: 200,
    });
  }
  async getByEmail(req: Request, res: Response) {
    const email = req.params.email as string;
    const customers = await this.service.getByEmail(email);
    sendResponse(res, {
      message: 'customer fetched successfully',
      data: customers,
      statusCode: 200,
    });
  }
}
