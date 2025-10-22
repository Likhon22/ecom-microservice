/* eslint-disable @typescript-eslint/no-explicit-any */
import { ZodError } from 'zod';
import handleZodError from '../../error/zodError.js';
import { Code, ConnectError } from '@bufbuild/connect';
import handleValidationError from '../../error/validationError.js';
import handleDuplicateError from '../../error/duplicateError.js';
import ApiError from '../../error/appError.js';

export function handleGrpcError(err: unknown) {
  if (err instanceof ZodError) {
    const simplified = handleZodError(err);
    return new ConnectError(
      simplified.message,
      Code.InvalidArgument,
      undefined,
      undefined,
      { errorSources: simplified.errorSources },
    );
  }
  if ((err as any)?.name === 'ValidationError') {
    const simplified = handleValidationError(err as any);
    return new ConnectError(
      simplified.message,
      Code.InvalidArgument,
      undefined,
      undefined,
      { errorSources: simplified.errorSources },
    );
  }
  if ((err as any)?.code === 11000) {
    const simplified = handleDuplicateError(err);
    return new ConnectError(
      simplified.message,
      Code.AlreadyExists,
      undefined,
      undefined,
      { errorSources: simplified.errorSources },
    );
  }
  if (err instanceof ApiError) {
    const grpcCode = mapHttpToGrpcCode(err.statusCode);
    return new ConnectError(err.message, grpcCode);
  }

  // Generic Error
  if (err instanceof Error) {
    return new ConnectError(err.message, Code.Internal);
  }
  return new ConnectError('Something went wrong', Code.Internal);
}

function mapHttpToGrpcCode(statusCode: number) {
  switch (statusCode) {
    case 400:
      return Code.InvalidArgument;
    case 401:
      return Code.Unauthenticated;
    case 403:
      return Code.PermissionDenied;
    case 404:
      return Code.NotFound;
    case 409:
      return Code.AlreadyExists;
    case 429:
      return Code.ResourceExhausted;
    case 500:
      return Code.Internal;
    case 501:
      return Code.Unimplemented;
    case 503:
      return Code.Unavailable;
    default:
      return Code.Unknown;
  }
}
