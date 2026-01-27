import { z, type ZodError } from 'zod';

export const validIpSchema = z.union([z.ipv4(), z.ipv6()]);

export const validStringSchema = z.string().trim().min(1);

export type FormatZodErrorOptions<T extends boolean = false> = {
  errorMessage?: string;
  preferSingleError?: T;
};

export type FormattedZodError<T extends boolean = false> = T extends true
  ? { message: string }
  : { message: string; details?: string[] };

export function formatZodError<T extends boolean = false>(
  error: ZodError,
  formatZodErrorOptions: FormatZodErrorOptions<T> = {}
): FormattedZodError<T> {
  const errors = error.issues.map(({ message, path }) => {
    const name = path.join('.');
    return name ? `${name} ${message}` : message;
  });

  let parsedMessage = formatZodErrorOptions.errorMessage ?? 'Invalid body';
  let parsedErrors: string[] | undefined = errors;

  if (formatZodErrorOptions.preferSingleError) {
    parsedMessage = errors[0];
    parsedErrors = undefined;
  }

  return {
    message: parsedMessage,
    ...(parsedErrors ? { details: parsedErrors } : {})
  };
}
