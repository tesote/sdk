/**
 * /status and /whoami response shapes.
 */

export interface StatusResponse {
  status: 'ok' | string;
  authenticated: boolean;
}

export type WhoamiClientType = 'workspace' | 'user';

export interface WhoamiClient {
  id: string;
  name: string;
  type: WhoamiClientType;
}

export interface WhoamiResponse {
  client: WhoamiClient;
}
