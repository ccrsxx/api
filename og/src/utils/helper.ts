import { type Request } from 'express';

export function getRemainingSecondsFromDate(date: Date): number {
  const now = new Date();
  return Math.max(0, Math.floor((date.getTime() - now.getTime()) / 1000));
}

export function getIpAddressFromRequest(req: Request): string {
  return req.ip ?? 'N/A';
}

export function getPublicUrlFromRequest(req: Request): string {
  const host = req.get('host');
  const protocol = req.protocol;
  return `${protocol}://${host}`;
}
