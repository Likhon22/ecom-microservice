import mongoose from 'mongoose';
import type {
  TErrorSources,
  TGenericErrorResponse,
} from '../types/error.type.js';

const handleValidationError = (
  err: mongoose.Error.ValidationError,
): TGenericErrorResponse => {
  const errorSources: TErrorSources = Object.values(err.errors).map(
    (val: mongoose.Error.ValidatorError | mongoose.Error.CastError) => {
      return {
        path: val.path,
        message: val.message,
      };
    },
  );

  return {
    statusCode: 400,

    message: 'Mongoose Validation Error',
    errorSources,
  };
};

export default handleValidationError;
