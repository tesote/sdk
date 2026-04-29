module TesoteSdk
  # Generic pagination helpers. Resource clients return native model objects;
  # these enumerators wrap repeated calls so callers can iterate the full set
  # without re-implementing cursor or offset arithmetic.
  module Pagination
    # Cursor pagination (transactions index): walks `pagination.after_id` until
    # has_more is false. Yields the raw page hash so callers can read the
    # response envelope as needed.
    class CursorEnumerator
      include Enumerable

      def initialize(start_query: {}, &fetch_page)
        raise ArgumentError, 'block (fetch_page) is required' unless block_given?

        @start_query = start_query.dup
        @fetch_page = fetch_page
      end

      def each
        return enum_for(:each) unless block_given?

        query = @start_query.dup
        loop do
          page = @fetch_page.call(query)
          yield page

          pagination = pagination_hash(page)
          break unless pagination['has_more']

          after_id = pagination['after_id']
          break if after_id.nil? || after_id.to_s.empty?

          query = query.merge(transactions_after_id: after_id)
        end
      end

      private

      def pagination_hash(page)
        return {} if page.nil?

        if page.is_a?(Hash)
          (page['pagination'] || page[:pagination] || {}).transform_keys(&:to_s)
        elsif page.respond_to?(:pagination)
          to_pagination_hash(page.pagination)
        else
          {}
        end
      end

      def to_pagination_hash(value)
        return {} if value.nil?
        return value.transform_keys(&:to_s) if value.is_a?(Hash)
        return value.to_h.transform_keys(&:to_s) if value.respond_to?(:to_h)

        {}
      end
    end

    # Offset pagination: walks until has_more is false, advancing `offset` by
    # `limit`. Yields each page hash.
    class OffsetEnumerator
      include Enumerable

      def initialize(start_query: {}, limit: 50, &fetch_page)
        raise ArgumentError, 'block (fetch_page) is required' unless block_given?

        @start_query = start_query.dup
        @limit = limit
        @fetch_page = fetch_page
      end

      def each
        return enum_for(:each) unless block_given?

        offset = (@start_query[:offset] || @start_query['offset'] || 0).to_i
        loop do
          query = @start_query.merge(limit: @limit, offset: offset)
          page = @fetch_page.call(query)
          yield page

          break unless page_has_more?(page)

          offset += @limit
        end
      end

      private

      def page_has_more?(page)
        return false if page.nil?
        return !!page['has_more'] if page.is_a?(Hash) && page.key?('has_more')
        return !!page[:has_more] if page.is_a?(Hash) && page.key?(:has_more)
        return !!page.has_more if page.respond_to?(:has_more)

        false
      end
    end
  end
end
