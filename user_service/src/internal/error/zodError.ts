import { ZodError, type ZodIssue } from 'zod';
import type {
  TErrorSources,
  TGenericErrorResponse,
} from '../types/error.type.js';

const handleZodError = (err: ZodError): TGenericErrorResponse => {
  const errorSources: TErrorSources = err.issues.map((issue: ZodIssue) => {
    return {
      path: issue?.path[issue.path.length - 1],
      message: issue.message,
    };
  });

  return {
    statusCode: 400,
    message: 'Validation Error',
    errorSources,
  };
};
export default handleZodError;
