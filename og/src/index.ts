import { createServer } from 'http';
import express from 'express';
import { logger } from './loaders/pino.ts';
import loaders from './loaders/index.ts';
import routes from './api/route.ts';
import errorHandler from './middleware/error.ts';

function main(): void {
  const app = express();
  const server = createServer(app);

  loaders(app, server);

  routes(app);

  errorHandler(app);

  const PORT = 4444;

  server.listen(PORT, () => {
    logger.info(`Server running on port ${PORT}`);
  });
}

main();
