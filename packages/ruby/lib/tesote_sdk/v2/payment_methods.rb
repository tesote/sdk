module TesoteSdk
  module V2
    class PaymentMethods
      def initialize(transport)
        @transport = transport
      end

      # GET /v2/payment_methods
      def list(query = {}, opts: {})
        body = @transport.request('GET', 'payment_methods', query: query, opts: opts)
        items_raw = (body && body['items']) || []
        Models::OffsetPage.new(
          items: Models::PaymentMethod.from_array(items_raw),
          limit: body && body['limit'],
          offset: body && body['offset'],
          has_more: body && body['has_more']
        )
      end

      # GET /v2/payment_methods/{id}
      def get(id, opts: {})
        raise ArgumentError, 'id is required' if id.nil? || id.to_s.empty?

        body = @transport.request('GET', "payment_methods/#{id}", opts: opts)
        Models::PaymentMethod.from_hash(body)
      end

      # POST /v2/payment_methods
      def create(payment_method:, opts: {})
        raise ArgumentError, 'payment_method is required' if payment_method.nil?

        payload = { payment_method: payment_method }
        body = @transport.request('POST', 'payment_methods', body: payload, opts: opts)
        Models::PaymentMethod.from_hash(body)
      end

      # PATCH /v2/payment_methods/{id}
      def update(id, payment_method:, opts: {})
        raise ArgumentError, 'id is required' if id.nil? || id.to_s.empty?
        raise ArgumentError, 'payment_method is required' if payment_method.nil?

        payload = { payment_method: payment_method }
        body = @transport.request('PATCH', "payment_methods/#{id}", body: payload, opts: opts)
        Models::PaymentMethod.from_hash(body)
      end

      # DELETE /v2/payment_methods/{id} → 204 No Content
      def delete(id, opts: {})
        raise ArgumentError, 'id is required' if id.nil? || id.to_s.empty?

        @transport.request('DELETE', "payment_methods/#{id}", opts: opts)
        nil
      end
    end
  end
end
