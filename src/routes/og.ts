import { Router, type Application } from 'express';
import { OgController } from '../controllers/og.ts';

export default (app: Application): void => {
  const router = Router();

  app.use('/og', router);

  router.get('/', OgController.getOg);
};
