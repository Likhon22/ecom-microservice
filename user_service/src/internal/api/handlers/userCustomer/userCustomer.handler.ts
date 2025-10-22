import type { Request, Response } from 'express';

import type { UserCustomerService } from '../../../service/userCustomer.service.js';
import type { UserCustomerDto } from '../../../domain/dtos/userCustmer.dto.js';
import sendResponse from '../../../utils/sendResponse.js';

export class UserCustomerHandler {
  private readonly service: UserCustomerService;
  constructor(service: UserCustomerService) {
    this.service = service;
  }

  async create(req: Request, res: Response) {
    const payload = req.body as UserCustomerDto;

    const user = await this.service.create(payload);
    sendResponse(res, {
      message: 'user created successfully',
      data: user,
      statusCode: 200,
    });
  }
}
