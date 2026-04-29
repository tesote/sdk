import type { StatusResponse, WhoamiResponse } from '../models/status.js';
import type { Transport } from '../transport.js';

export class V1StatusClient {
  constructor(private readonly transport: Transport) {}

  /** GET /status — auth not required. */
  async status(): Promise<StatusResponse> {
    const res = await this.transport.request<StatusResponse>({
      method: 'GET',
      path: '/status',
    });
    return res.data;
  }

  /** GET /whoami — auth required. */
  async whoami(): Promise<WhoamiResponse> {
    const res = await this.transport.request<WhoamiResponse>({
      method: 'GET',
      path: '/whoami',
    });
    return res.data;
  }
}
