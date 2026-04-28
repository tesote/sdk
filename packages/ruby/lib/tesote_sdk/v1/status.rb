module TesoteSdk
  module V1
    class Status
      def initialize(transport)
        @transport = transport
      end

      def status(opts: {})
        raise NotImplementedError, 'V1::Status#status not wired in 0.1.0'
      end

      def whoami(opts: {})
        raise NotImplementedError, 'V1::Status#whoami not wired in 0.1.0'
      end
    end
  end
end
