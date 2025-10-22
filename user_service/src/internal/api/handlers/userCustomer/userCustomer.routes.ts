import { Router } from 'express';
import type { UserCustomerHandler } from './userCustomer.handler.js';
import type { Middleware } from '../../middleware/index.js';
import { TCustomerSchema } from './userCustomer.validation.js';
import { catchAsync } from '../../../utils/catchAsync.js';

export function userCustomerRoutes(
  handler: UserCustomerHandler,
  mw: Middleware,
): Router {
  const router = Router();
  router.post(
    '/',
    mw.validate(TCustomerSchema),
    catchAsync(handler.create.bind(handler)),
  );
  router.get('/', catchAsync(handler.get.bind(handler)));

  return router;
}
