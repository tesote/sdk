import type { StatusResponse, WhoamiResponse } from '../models/status.js';
import type { Transport } from '../transport.js';

export class V2StatusClient {
  constructor(private readonly transport: Transport) {}

  /** GET /v2/status — auth not required. */
  async status(): Promise<StatusResponse> {
    const res = await this.transport.request<StatusResponse>({
      method: 'GET',
      path: '/v2/status',
    });
    return res.data;
  }

  /** GET /v2/whoami — auth required. */
  async whoami(): Promise<WhoamiResponse> {
    const res = await this.transport.request<WhoamiResponse>({
      method: 'GET',
      path: '/v2/whoami',
    });
    return res.data;
  }
}
