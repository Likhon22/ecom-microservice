import mongoose from 'mongoose';
import type {
  TErrorSources,
  TGenericErrorResponse,
} from '../types/error.type.js';

const handleCastError = (
  err: mongoose.Error.CastError,
): TGenericErrorResponse => {
  const errorSources: TErrorSources = [
    {
      path: err.path,
      message: err.message,
    },
  ];
  return {
    statusCode: 400,
    message: 'Invalid Id',
    errorSources: errorSources,
  };
};

export default handleCastError;
