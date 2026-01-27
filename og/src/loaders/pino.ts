import { pino as Pino, type LoggerOptions } from 'pino';
import { pinoHttp as PinoHttp, type Options } from 'pino-http';
import { appConfig } from '../config/config.ts';
import type { Application } from 'express';

type CombinedLoggerOptions = {
  pinoOptions: LoggerOptions;
  pinoHttpOptions?: Options;
};

const developmentLoggerOptions: CombinedLoggerOptions = {
  pinoOptions: {
    transport: {
      target: 'pino-pretty'
    }
  },
  pinoHttpOptions: {
    autoLogging: false
  }
};

const productionLoggerOptions: CombinedLoggerOptions = {
  pinoOptions: {
    level: 'trace',
    formatters: {
      level(label) {
        return { severity: label };
      }
    },
    messageKey: 'message'
  }
};

const { pinoOptions, pinoHttpOptions } = appConfig.isDevelopment
  ? developmentLoggerOptions
  : productionLoggerOptions;

export const logger = Pino(pinoOptions);

const pinoHttp = PinoHttp({
  ...pinoHttpOptions,
  logger
});

export default (app: Application): void => {
  app.use(pinoHttp);
};
