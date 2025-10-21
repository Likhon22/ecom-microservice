import type { Response } from 'express';

type TResponse<T> = {
  message: string;
  statusCode: number;
  data: T;
};

function sendResponse<T>(res: Response, responseData: TResponse<T>) {
  res.status(responseData.statusCode).json({
    message: responseData.message,
    data: responseData.data,
  });
}
export default sendResponse;
