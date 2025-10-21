import type { Request, Response } from 'express';

import type { UserCustomerService } from '../../../service/userCustomer.service.js';
import type { UserCustomerDto } from '../../../domain/dtos/userCustmer.dto.js';

export class UserCustomerHandler {
  private readonly service: UserCustomerService;
  constructor(service: UserCustomerService) {
    this.service = service;
  }

  async create(req: Request, res: Response) {
    try {
      const payload = req.body as UserCustomerDto;
      console.log(payload);

      const user = await this.service.create(payload);
      res.send(user);
    } catch (err) {
      console.log(err);
    }
  }
}
