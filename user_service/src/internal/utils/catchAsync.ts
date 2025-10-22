import type { NextFunction, Request, Response } from 'express';
import type { RequestHandler } from 'express-serve-static-core';

export const catchAsync = (fn: RequestHandler) => {
  return (req: Request, res: Response, next: NextFunction) => {
    Promise.resolve(fn(req, res, next)).catch(err => next(err));
  };
};
