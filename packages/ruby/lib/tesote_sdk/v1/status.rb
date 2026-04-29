module TesoteSdk
  module V1
    # GET /status, GET /whoami — both live under the API root, NOT under /v1.
    # We bypass the transport's version_segment by using opts[:extra_headers] is
    # not enough; we need an absolute path. The transport joins
    # base_url + version_segment + path, so we send a leading double-slash
    # path? No — instead we use an unversioned subclient via raw_request.
    class Status
      def initialize(transport)
        @transport = transport
      end

      # GET /status (no auth required, but transport always sends it — server
      # ignores when not required).
      def status(opts: {})
        body = @transport.request_unversioned('GET', 'status', opts: opts)
        Models::StatusResult.from_hash(body)
      end

      # GET /whoami
      def whoami(opts: {})
        body = @transport.request_unversioned('GET', 'whoami', opts: opts)
        Models::Whoami.from_hash(body)
      end
    end
  end
end
