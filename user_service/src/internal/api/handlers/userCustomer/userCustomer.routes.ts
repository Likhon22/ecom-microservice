import { Router } from 'express';
import type { UserCustomerHandler } from './userCustomer.handler.js';
import type { Middleware } from '../../middleware/index.js';
import { TCustomerSchema } from './userCustomer.validation.js';

export function userCustomerRoutes(
  handler: UserCustomerHandler,
  mw: Middleware,
): Router {
  const router = Router();
  router.post('/', mw.validate(TCustomerSchema), handler.create.bind(handler));

  return router;
}
