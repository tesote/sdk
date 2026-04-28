require 'openssl'

module TesoteSdk
  module V3
    # Stateless helper. Confirms the platform's signature scheme is HMAC-SHA256
    # of the raw request body keyed by the webhook secret. Confirm before
    # depending in production — see docs/architecture/resources.md.
    module WebhooksSignature
      module_function

      def verify_webhook_signature(body:, signature_header:, secret:)
        if body.nil? || signature_header.nil? || signature_header.empty? || secret.nil? || secret.empty?
          raise ArgumentError, 'body, signature_header, and secret are required'
        end

        expected = compute(body, secret)
        return true if secure_compare(expected, signature_header)

        raise TesoteSdk::Error.new('webhook signature mismatch',
                                   error_code: 'WEBHOOK_SIGNATURE_MISMATCH')
      end

      def compute(body, secret)
        OpenSSL::HMAC.hexdigest('SHA256', secret, body)
      end

      def secure_compare(a, b)
        return false unless a.bytesize == b.bytesize

        OpenSSL.fixed_length_secure_compare(a, b)
      end
    end

    # Re-export at module level for the README example: TesoteSdk::V3.verify_webhook_signature(...)
    def self.verify_webhook_signature(body:, signature_header:, secret:)
      WebhooksSignature.verify_webhook_signature(body: body, signature_header: signature_header, secret: secret)
    end
  end
end
