/* eslint-disable no-console */

import dotenv from 'dotenv/config';
import { writeFile } from 'fs/promises';
import { execSync } from 'child_process';
import { transpile } from 'postman2openapi';
import type { OpenAPIV3 } from 'openapi-types';

async function main(): Promise<void> {
  console.log('Generating OpenAPI documentation...');

  const postmanEnvSchema = {
    POSTMAN_API_KEY: process.env.POSTMAN_API_KEY,
    POSTMAN_COLLECTION_ID: process.env.POSTMAN_COLLECTION_ID
  };

  const isPostmanEnvMissing =
    !postmanEnvSchema.POSTMAN_API_KEY ||
    !postmanEnvSchema.POSTMAN_COLLECTION_ID;

  if (isPostmanEnvMissing) {
    throw new Error(
      'Missing Postman environment variables. Please set POSTMAN_API_KEY and POSTMAN_COLLECTION_ID.'
    );
  }

  const postmanCollectionResponse = await fetch(
    `https://api.getpostman.com/collections/${postmanEnvSchema.POSTMAN_COLLECTION_ID}`,
    {
      headers: {
        'x-api-key': postmanEnvSchema.POSTMAN_API_KEY as string
      }
    }
  );

  if (!postmanCollectionResponse.ok) {
    throw new Error(
      `Failed to fetch Postman collection: ${postmanCollectionResponse.statusText}`
    );
  }

  type PostmanCollection = {
    collection: {
      info: Record<string, unknown>;
      item: Record<string, unknown>[];
    };
  };

  const data = (await postmanCollectionResponse.json()) as PostmanCollection;

  const openapi = transpile(data.collection) as OpenAPIV3.Document;

  const stringifiedOpenapi = JSON.stringify(openapi);

  const outputPath = '../internal/features/docs/openapi.json';

  await writeFile(outputPath, stringifiedOpenapi);

  console.log('OpenAPI documentation generated successfully.');

  const skipCommit = process.argv.includes('--no-commit');

  if (skipCommit) {
    console.log(
      'Skipping commit due to --no-commit flag. Please commit manually.'
    );
    return;
  }

  try {
    execSyncWithOutput(`git diff --quiet --exit-code ${outputPath}`);

    console.log('No changes detected in the OpenAPI spec. Nothing to commit.');

    return;
  } catch {
    console.log(
      'Changes detected in the OpenAPI spec. Proceeding to commit...'
    );
  }

  execSyncWithOutput(`git add ${outputPath}`);

  console.log('Committing changes...');

  execSyncWithOutput('git commit -m "docs(api): update OpenAPI spec"');

  console.log('Successfully committed API documentation changes.');
}

function execSyncWithOutput(command: string): void {
  execSync(command, { stdio: 'inherit' });
}

void main().catch((error) => {
  console.error(error);
  process.exit(1);
});
