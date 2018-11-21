module GromNative
    # URL request object, allowing the user to build a URL to make a request to an API.
    #
    # @since 0.7.5
    #
    # @attr_reader [String] base_url the endpoint for our API which we will build our requests on. (expected: http://example.com - without the trailing slash).
    # @attr_reader [Hash] headers the headers being sent in the request.
    class UrlRequest < Parliament::Request::BaseRequest
      # Creates a new instance of Parliament::Request::UrlRequest.
      #
      # @see Parliament::Request::BaseRequest#initialize.
      #
      # @param [String] base_url the base url of our api. (expected: http://example.com - without the trailing slash).
      # @param [Hash] headers the headers being sent in the request.
      # @param [Parliament::Builder] builder the builder to use in order to build a response.
      # @param [Module] decorators the decorator module to use in order to provide possible alias methods for any objects created by the builder.
      # @example Passing headers
      #
      # request = Parliament::Request::UrlRequest.new(base_url: 'http://example.com', headers: { 'Access-Token' => '12345678' })
      # This will create a request with the Access-Token set to 12345678.
      def initialize(base_url: nil, headers: nil, builder: nil, decorators: nil)
        @endpoint_parts = []
        @base_url     ||= ENV['PARLIAMENT_BASE_URL']
        @headers        = headers || self.class.headers || {}
        @builder        = builder || Parliament::Builder::BaseResponseBuilder
        @decorators     = decorators
        @query_params   = {}
      end

      # Makes an HTTP GET request and process results into a response.
      #
      # @example HTTP GET request
      #   request = Parliament::Request::BaseRequest.new(base_url: 'http://example.com/people/123'
      #
      #   # url: http://example.com/people/123
      #
      #   response = request.get #=> #<Parliament::Response::BaseResponse ...>
      #
      # @example HTTP GET request with URI encoded form values
      #   request = Parliament::Request.new(base_url: 'http://example.com/people/current')
      #
      #   # url: http://example.com/people/current?limit=10&page=4&lang=en-gb
      #
      #   response = request.get({ limit: 10, page: 4, lang: 'en-gb' }) #=> #<Parliament::Response::BaseResponse ...>
      #
      # @raise [Parliament::ServerError] when the server responds with a 5xx status code.
      # @raise [Parliament::ClientError] when the server responds with a 4xx status code.
      # @raise [Parliament::NoContentResponseError] when the response body is empty.
      #
      # @param [Hash] params (optional) additional URI encoded form values to be added to the URI.
      #
      # @return [Parliament::Response::BaseResponse] a Parliament::Response::BaseResponse object containing all of the data returned from the URL.
      def get(params: {}, filter: [])
        uri = URI.parse(query_url)

        temp_params = {}

        if uri.query
          # Returns [ ["key", "value"], ["key", "value"] ]
          key_value_array = URI.decode_www_form(endpoint.query)
          key_value_array.map! { |key_value_pair| [ key_value_pair[0].to_sym, key_value_pair[1] ] }
          temp_params = key_value_array.to_h
        end

        temp_params = temp_params.merge(params)

        uri.query = temp_params

        GromNative.fetch( uri: uri.to_s, headers: headers, filter: filter, decorators: @decorators )
      end

      # Overrides ruby's method_missing to allow creation of URLs through method calls.
      #
      # @example Adding a simple URL part
      #   request = Parliament::Request::UrlRequest.new(base_url: 'http://example.com')
      #
      #   # url: http://example.com/people
      #   request.people
      #
      # @example Adding a simple URL part with parameters
      #   request = Parliament::Request::UrlRequest.new(base_url: 'http://example.com')
      #
      #   # url: http://example.com/people/123456
      #   request.people('123456')
      #
      # @example Chaining URL parts and using hyphens
      #   request = Parliament::Request::UrlRequest.new(base_url: 'http://example.com')
      #
      #   # url: http://example.com/people/123456/foo/bar/hello-world/7890
      #   request.people('123456').foo.bar('hello-world', '7890')
      #
      # @param [Symbol] method the 'method' (url part) we are processing.
      # @param [Array<Object>] params parameters passed to the specified method (url part).
      # @param [Block] block additional block (kept for compatibility with method_missing API).
      #
      # @return [Parliament::Request::UrlRequest] self (this is to allow method chaining).
      def method_missing(method, *params, &block)
        @endpoint_parts << method.to_s

        @endpoint_parts << params
        @endpoint_parts = @endpoint_parts.flatten!

        block&.call

        self || super
      end

      # This class always responds to method calls, even those missing. Therefore, respond_to_missing? always returns true.
      #
      # @return [Boolean] always returns true.
      def respond_to_missing?(_, _ = false)
        true # responds to everything, always
      end

      def query_url
        uri_string = [@base_url, @endpoint_parts].join('/').gsub(' ', '%20')

        uri = URI.parse(uri_string)
        uri.query = URI.encode_www_form(@query_params) unless @query_params.empty?

        uri.to_s
      end

      # @return [Parliament::Request::UrlRequest] self (this is to allow method chaining).
      def set_url_params(params)
        @query_params = @query_params.merge(params)

        self
      end

      def default_headers
        { 'Accept' => ['*/*', 'application/n-triples'] }
      end

      def headers
        default_headers.merge(@headers)
      end

      class << self
        attr_accessor :base_url, :headers
      end
    end
  end
end