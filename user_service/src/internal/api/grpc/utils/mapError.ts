/* eslint-disable @typescript-eslint/no-explicit-any */
import { status, Metadata, type ServiceError } from '@grpc/grpc-js';
import { ZodError } from 'zod';
import ApiError from '../../../error/appError.js';
import handleZodError from '../../../error/zodError.js';
import handleValidationError from '../../../error/validationError.js';
import handleDuplicateError from '../../../error/duplicateError.js';
import handleCastError from '../../../error/castError.js';

export function mapGrpcError(err: unknown): ServiceError {
  if (err instanceof ZodError) {
    const simplified = handleZodError(err);
    return buildError(simplified.message, status.INVALID_ARGUMENT);
  }
  if ((err as any)?.name === 'ValidationError') {
    const simplified = handleValidationError(err as any);
    return buildError(simplified.message, status.INVALID_ARGUMENT);
  }
  if ((err as any)?.name === 'CastError') {
    const simplified = handleCastError(err as any);
    return buildError(simplified.message, status.INVALID_ARGUMENT);
  }
  if ((err as any)?.code === 11000) {
    const simplified = handleDuplicateError(err as any);
    return buildError(simplified.message, status.ALREADY_EXISTS);
  }
  if (err instanceof ApiError) {
    return buildError(err.message, httpToGrpc(err.statusCode));
  }
  if (err instanceof Error) {
    return buildError(err.message, status.UNKNOWN);
  }
  return buildError('Something went wrong', status.UNKNOWN);
}

function buildError(message: string, code: status): ServiceError {
  const error = new Error(message) as ServiceError;
  error.code = code;
  error.details = message;
  error.metadata = new Metadata();
  return error;
}

function httpToGrpc(statusCode: number): status {
  switch (statusCode) {
    case 400:
      return status.INVALID_ARGUMENT;
    case 401:
      return status.UNAUTHENTICATED;
    case 403:
      return status.PERMISSION_DENIED;
    case 404:
      return status.NOT_FOUND;
    case 409:
      return status.ALREADY_EXISTS;
    case 429:
      return status.RESOURCE_EXHAUSTED;
    case 500:
      return status.INTERNAL;
    case 501:
      return status.UNIMPLEMENTED;
    case 503:
      return status.UNAVAILABLE;
    default:
      return status.UNKNOWN;
  }
}
