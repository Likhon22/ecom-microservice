import type { Application } from 'express';
import type { UserCustomerHandler } from '../handlers/userCustomer/userCustomer.handler.js';
import { userCustomerRoutes } from '../handlers/userCustomer/userCustomer.routes.js';
import type { Middleware } from '../middleware/index.js';

export function registerRoutes(
  app: Application,
  userHandler: UserCustomerHandler,
  mw: Middleware,
) {
  // attach all routes under versioned API prefix
  app.use('/api/v1/usercustomers', userCustomerRoutes(userHandler, mw));
}
