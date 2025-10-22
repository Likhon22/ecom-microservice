import globalErrorHandler from './globalErrorHanlder.js';
import notFound from './notFound.js';
import validateRequest from './validateRequest.js';

export class Middleware {
  validate = validateRequest;
  noRoutesFound = notFound;
  globalErrorHandler = globalErrorHandler;
}
