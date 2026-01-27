import { randomUUID } from 'crypto';
import { HttpError, FatalError, AppError, errorAs } from '../api/error.ts';
import { logger } from '../loaders/pino.ts';
import { getIpAddressFromRequest } from '../utils/helper.ts';
import type { ApiResponse } from '../types/api.ts';
import type { Request, Response, Application, NextFunction } from 'express';
import type { ApiLogContext } from '../types/log.ts';

export default (app: Application): void => {
  app.use(notFoundHandler);
  app.use(errorHandler);
};

function notFoundHandler(
  req: Request,
  _res: Response,
  next: NextFunction
): void {
  const notFoundError = new HttpError(404, {
    message: `Route not found - ${req.originalUrl}`
  });

  next(notFoundError);
}

function errorHandler(
  err: Error,
  req: Request,
  res: Response<ApiResponse>,
  _next: NextFunction
): void {
  const errorId = randomUUID();

  type ApiLogContext<T = Record<string, unknown>> = {
    context: T;
  };

  type LogContext = {
    ip: string;
    url: string;
    error: Error;
    method: string;
    errorId: string;
  };

  const logContext: ApiLogContext<LogContext> = {
    context: {
      ip: getIpAddressFromRequest(req),
      url: req.originalUrl,
      error: err,
      method: req.method,
      errorId: errorId
    }
  };

  logger.debug(logContext, 'Error handler invoked');

  const fatalError = errorAs(err, FatalError);

  if (fatalError) {
    logger.fatal(logContext, `Handled fatal error - ${err.message}`);

    process.exit(1);
  }

  const httpError = errorAs(err, HttpError);

  if (httpError) {
    logger.info(logContext, `Handled expected error - ${err.message}`);

    res.status(httpError.statusCode).json({
      error: { id: errorId, message: err.message, details: httpError.details }
    });

    return;
  }

  const appError = errorAs(err, AppError);

  if (appError) {
    logger.warn(logContext, `Handled application error - ${err.message}`);

    res.status(400).json({
      error: { id: errorId, message: err.message, details: [] }
    });

    return;
  }

  // Below this point, we treat the error as an unexpected error

  // Any unhandled error
  if (err instanceof Error) {
    logger.error(logContext, `Handled unexpected error - ${err.message}`);

    res.status(500).json({
      error: {
        id: errorId,
        message: 'An internal server error occurred.',
        details: []
      }
    });

    return;
  }

  // Fallback for non-Error throwables (e.g., throw 'string')

  logger.error(logContext, 'Handled unknown error');

  res.status(500).json({
    error: {
      id: errorId,
      message: 'An internal server error occurred.',
      details: []
    }
  });
}
