type AppErrorOptions = ErrorOptions & {
  message: string;
};

type HttpErrorOptions = AppErrorOptions & {
  details?: string[];
};

export class HttpError extends Error {
  public statusCode: number;
  public details: string[];

  public constructor(statusCode: number, options: HttpErrorOptions) {
    const { message, details = [], cause } = options;

    super(message, { cause });

    this.details = details;
    this.statusCode = statusCode;
  }
}

export class AppError extends Error {
  public constructor(options: AppErrorOptions) {
    const { message, cause } = options;

    super(message, { cause });
  }
}

export class FatalError extends Error {
  public constructor(options: AppErrorOptions) {
    const { message, cause } = options;

    super(message, { cause });
  }
}

export function errorAs<T extends Error>(
  err: unknown,
  targetErr: new (...args: never[]) => T
): T | null {
  let currentErr = err;

  while (currentErr instanceof Error) {
    if (currentErr instanceof targetErr) return currentErr;

    currentErr = currentErr.cause;
  }

  return null;
}
