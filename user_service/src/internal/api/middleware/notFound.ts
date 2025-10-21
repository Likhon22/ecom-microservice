/* eslint-disable @typescript-eslint/no-unused-vars */
import type { NextFunction, Request, Response } from 'express';
import type {
  TGenericErrorResponse,
  TGenericUserErrorResponse,
} from '../../types/error.type.js';

const notFound = (req: Request, res: Response, next: NextFunction): void => {
  const errorResponse: TGenericErrorResponse = {
    statusCode: 404,
    message: 'API is not found',
    errorSources: [{ path: req.originalUrl, message: 'API is not found' }],
  };
  console.log(errorResponse);
  const userErrorResponse: TGenericUserErrorResponse = {
    statusCode: 404,
    message: 'API is not found',
  };

  res.status(errorResponse.statusCode).json(userErrorResponse);
};

export default notFound;
