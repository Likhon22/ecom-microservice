/* eslint-disable @typescript-eslint/no-explicit-any */

import type {
  TErrorSources,
  TGenericErrorResponse,
} from '../types/error.type.js';

const handleDuplicateError = (err: any): TGenericErrorResponse => {
  const message = `Duplicate field value`;

  const regex = /{ name: "(.*?)" }/;
  const match = err?.errmsg.match(regex);
  const errorMessage = match ? match[1] : '';

  const errorSources: TErrorSources = [
    {
      path: 'name',
      message: `${errorMessage} is already exists`,
    },
  ];
  return {
    statusCode: 400,
    message,
    errorSources,
  };
};
export default handleDuplicateError;
