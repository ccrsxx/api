import og from '../features/og/route.ts';
import type { Application } from 'express';

export default (app: Application): void => {
  og(app);
};
