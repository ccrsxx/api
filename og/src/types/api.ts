export type SuccessApiResponse<T> = {
  data: T | null;
};

export type ErrorApiResponse = {
  error: {
    id: string;
    message: string;
    details: string[];
  };
};

export type ApiResponse<T = unknown> = SuccessApiResponse<T> | ErrorApiResponse;
