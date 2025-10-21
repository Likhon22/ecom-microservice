import type { NextFunction, Request, Response } from 'express';
import type { ZodTypeAny } from 'zod';

const validateRequest = (schema: ZodTypeAny) => {
  return async (req: Request, res: Response, next: NextFunction) => {
    try {
      await schema.parseAsync({
        body: req.body,
        cookies: req.cookies,
      });
      next();
    } catch (err) {
      next(err);
    }
  };
};

export default validateRequest;
