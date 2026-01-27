import pino from './pino.ts';
import type { Application } from 'express';
import type { Server } from 'http';

export default (app: Application, _server: Server): void => {
  pino(app);
};
