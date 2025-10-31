import z, { ZodError } from 'zod';
import type {
  TErrorSources,
  TGenericErrorResponse,
} from '../types/error.type.js';

const handleZodError = (err: ZodError): TGenericErrorResponse => {
  let zodMsg;
  const errorSources: TErrorSources = err.issues.map((issue: z.ZodIssue) => {
    zodMsg = issue.message;
    return {
      path: issue?.path[issue.path.length - 1] ?? '',
      message: issue.message,
    };
  });

  return {
    statusCode: 400,
    message: zodMsg || 'Validation error',
    errorSources,
  };
};
export default handleZodError;
