import type { Transport } from '../transport.js';

export class V3WebhooksClient {
  constructor(private readonly _transport: Transport) {}

  async list(_params: Record<string, unknown> = {}): Promise<unknown> {
    throw new Error('not implemented');
  }

  async get(_id: string): Promise<unknown> {
    throw new Error('not implemented');
  }

  async create(
    _body: Record<string, unknown>,
    _opts: { idempotencyKey?: string } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }

  async update(
    _id: string,
    _body: Record<string, unknown>,
    _opts: { idempotencyKey?: string } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }

  async delete(_id: string, _opts: { idempotencyKey?: string } = {}): Promise<unknown> {
    throw new Error('not implemented');
  }
}

export interface VerifyWebhookSignatureInput {
  body: string;
  signatureHeader: string;
  secret: string;
}

/**
 * Verify a webhook payload signature.
 *
 * The platform's signature scheme is not yet finalized — this stub throws.
 * See docs/architecture/resources.md ("Webhook signature verification").
 */
export function verifyWebhookSignature(_input: VerifyWebhookSignatureInput): void {
  throw new Error('not implemented');
}
