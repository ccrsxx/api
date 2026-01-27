import { OgService } from './service.ts';
import type { Request, Response } from 'express';

async function getOg(req: Request, res: Response<Buffer>): Promise<void> {
  const og = await OgService.getOg(req.query);

  res.contentType('image/png').send(og);
}

export const OgController = {
  getOg
};
