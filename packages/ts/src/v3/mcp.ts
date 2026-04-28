import type { Transport } from '../transport.js';

export class V3McpClient {
  constructor(private readonly _transport: Transport) {}

  /** POST /v3/mcp — pass-through; SDK exposes raw call, not a parsed model. */
  async handle(
    _body: Record<string, unknown>,
    _opts: { idempotencyKey?: string } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }
}
