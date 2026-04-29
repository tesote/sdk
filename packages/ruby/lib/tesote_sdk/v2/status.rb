module TesoteSdk
  module V2
    class Status
      def initialize(transport)
        @transport = transport
      end

      # GET /v2/status — note this lives at /api/v2/status (not /api/status).
      def status(opts: {})
        body = @transport.request('GET', 'status', opts: opts)
        Models::StatusResult.from_hash(body)
      end

      # GET /v2/whoami
      def whoami(opts: {})
        body = @transport.request('GET', 'whoami', opts: opts)
        Models::Whoami.from_hash(body)
      end
    end
  end
end
