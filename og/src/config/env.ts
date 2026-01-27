import { access } from 'fs/promises';
import dotenv from 'dotenv';
import { z } from 'zod';
import { logger } from '../loaders/pino.ts';
import { appConfig } from '../config/config.ts';
import { validStringSchema } from '../utils/validation.ts';

const envSchema = z.object({
  OWNER_NAME: validStringSchema,
  FRONTEND_DOMAIN: validStringSchema,
  CLOUDFLARE_CDN_URL: validStringSchema
});

type EnvSchema = z.infer<typeof envSchema>;

function validateEnv(): EnvSchema {
  const PORT = process.env.PORT ?? process.env.HOST_PORT;

  const mergedEnv = {
    ...process.env,
    PORT
  };

  let { data, error } = envSchema.safeParse(mergedEnv);

  const runningOnCi = process.env.CI === 'true';

  if (runningOnCi) {
    data = process.env as EnvSchema;
  }

  const shouldThrowError = error && !runningOnCi;

  if (shouldThrowError) {
    throw new Error(`Environment validation error: ${error.message}`);
  }

  return data as EnvSchema;
}

async function loadEnv(): Promise<void> {
  let envPath = '.env';

  if (appConfig.isDevelopment) {
    const isLocalEnvExists = await access('.env.local')
      .then(() => true)
      .catch(() => false);

    if (isLocalEnvExists) {
      envPath = '.env.local';
    }

    logger.info(`Loading environment variables from ${envPath}`);
  } else {
    logger.info('Loading environment variables from .env or process.env');
  }

  dotenv.config({ path: envPath, quiet: true });
}

await loadEnv();

export const appEnv = validateEnv();
