import type { ErrorRequestHandler } from 'express';
import type { TErrorSources } from '../../types/error.type.js';
import { ZodError } from 'zod';
import handleZodError from '../../error/zodError.js';
import handleValidationError from '../../error/validationError.js';
import handleCastError from '../../error/castError.js';
import handleDuplicateError from '../../error/duplicateError.js';
import ApiError from '../../error/appError.js';
import config from '../../config/index.js';

export const globalErrorHandler: ErrorRequestHandler = (
  err,
  req,
  res,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  next,
) => {
  let statusCode = err.statusCode || 500;
  let message = err.message || 'Something went wrong';

  let errorSources: TErrorSources = [
    { path: '', message: 'Something went wrong' },
  ];

  if (err instanceof ZodError) {
    const simpliFiedError = handleZodError(err);
    statusCode = simpliFiedError?.statusCode;
    message = simpliFiedError?.message;
    errorSources = simpliFiedError?.errorSources;
  } else if (err?.name === 'ValidationError') {
    const simpliFiedError = handleValidationError(err);
    statusCode = simpliFiedError?.statusCode;
    message = simpliFiedError?.message;
    errorSources = simpliFiedError?.errorSources;
  } else if (err?.name === 'CastError') {
    const simpliFiedError = handleCastError(err);
    statusCode = simpliFiedError?.statusCode;
    message = simpliFiedError?.message;
    errorSources = simpliFiedError?.errorSources;
  } else if (err?.code === 11000) {
    const simpliFiedError = handleDuplicateError(err);
    statusCode = simpliFiedError?.statusCode;
    message = simpliFiedError?.message;
    errorSources = simpliFiedError?.errorSources;
  } else if (err instanceof ApiError) {
    statusCode = err.statusCode;
    message = err.message;
    errorSources = [
      {
        path: '',
        message: err.message,
      },
    ];
  } else if (err instanceof Error) {
    message = err.message;
    errorSources = [
      {
        path: '',
        message: err.message,
      },
    ];
  }
  res.status(statusCode).json({
    success: false,
    message,
    errorSources,
    stack:
      config.node_env === 'development' && err.stack ? err.stack : undefined,
  });
};

export default globalErrorHandler;
