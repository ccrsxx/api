import { getPublicUrlFromRequest } from '../utils/helper.ts';
import type { Request, Response } from 'express';
import type { ApiResponse } from '../utils/types/api.ts';

function ping(req: Request, res: Response<ApiResponse>): void {
  res.status(200).json({
    data: {
      message: 'Welcome! The API is up and running',
      documentationUrl: `${getPublicUrlFromRequest(req)}/docs`
    }
  });
}

export const RootController = {
  ping
};
