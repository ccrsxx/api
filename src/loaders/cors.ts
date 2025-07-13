import cors, { type CorsOptions } from 'cors';
import { appEnv } from '../lib/env.js';
import type { Application } from 'express';

const ALLOWED_ORIGINS = appEnv.VALID_ORIGINS.split(',');

export const corsOptions: CorsOptions = {
  origin: ALLOWED_ORIGINS,
  methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE'],
  allowedHeaders: ['Authorization', 'Content-Type'],
  credentials: true
};

export default (app: Application): void => {
  app.use(cors(corsOptions));
};
