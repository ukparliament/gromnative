module GromNative
  # A Ruby object populated with n-triple data.
  #
  # @since 0.1.0
  # @attr_reader [Array] statements an array of n-triple statements.
  class Node
    BLANK = 'blank_node'.freeze

    attr_reader :statements

    # @param [Array] statements an array of n-triple statements.
    def initialize(statements, decorators = nil)
      @statements = statements

      populate(decorators)
    end

    # Allows the user to access instance variables as methods or raise an error if the variable is not defined.
    #
    # @param [Symbol] method name of method.
    # @param [Array] *params extra arguments to pass to super.
    # @param [Block] &block block to pass to super.
    # @example Accessing instance variables populated from statements
    #   statements = [
    #      RDF::Statement.new(RDF::URI.new('http://example.com/123'), RDF.type, 'Person'),
    #      RDF::Statement.new(RDF::URI.new('http://example.com/123'), RDF::URI.new('http://example.com/forename'), 'Jane'),
    #      RDF::Statement.new(RDF::URI.new('http://example.com/123'), RDF::URI.new('http://example.com/surname'), 'Smith')
    #   ]
    #
    #   node = Grom::Node.new(statements)
    #
    #   node.forename #=> 'Jane'
    #
    # @example Accessing instance variables created on the fly
    #   statements = [RDF::Statement.new(RDF::URI.new('http://example.com/123'), RDF.type, 'Person')]
    #
    #   node = Grom::Node.new(statements)
    #   node.instance_variable_set('@foo', 'bar')
    #
    #   node.foo #=> 'bar'
    #
    # @raise [NoMethodError] raises error if the method does not exist.
    def method_missing(method, *params, &block)
      instance_variable_get("@#{method}".to_sym) || super
    end

    # Allows the user to check if a Grom::Node responds to an instance variable
    #
    # @param [Symbol] method name of method.
    # @param [Boolean] include_all indicates whether to include private and protected methods (defaults to false).
    # @example Using respond_to?
    #
    #   statements = [
    #      RDF::Statement.new(RDF::URI.new('http://example.com/123'), RDF.type, 'Person'),
    #      RDF::Statement.new(RDF::URI.new('http://example.com/123'), RDF::URI.new('http://example.com/forename'), 'Jane'),
    #      RDF::Statement.new(RDF::URI.new('http://example.com/123'), RDF::URI.new('http://example.com/surname'), 'Smith')
    #   ]

    #   node = Grom::Node.new(statements)
    #
    #   node.respond_to?(:forename) #=> 'Jane'
    #   node.respond_to?(:foo) #=> false
    def respond_to_missing?(method, include_all = false)
      instance_variable_get("@#{method}".to_sym) || super
    end

    # Checks if Grom::Node is a blank node
    #
    # @return [Boolean] a boolean depending on whether or not the Grom::Node is a blank node
    def blank?
      @statements.first.fetch('subject', '').match(%r(^_:))
    end

    private

    def set_graph_id
      graph_id = Grom::Helper.get_id(@statements.first['subject'])
      instance_variable_set('@graph_id'.to_sym, graph_id)
    end

    def populate(decorators)
      set_graph_id
      @statements.each do |statement|
        predicate = Grom::Helper.get_id(statement['predicate']).to_sym
        object = RDF::NTriples::Reader.parse_object(statement['object'])

        object = if object.is_a? RDF::URI
                   object.to_s
                 else
                   object.object
                 end

        instance_variable = instance_variable_get("@#{predicate}")

        if instance_variable
          instance_variable = instance_variable.is_a?(Array) ? instance_variable.flatten : [instance_variable]
          instance_variable << object
          instance_variable_set("@#{predicate}", instance_variable)
        else
          instance_variable_set("@#{predicate}", object)
        end

        decorators&.decorate_with_type(self, object) if statement['predicate'] == RDF.type && decorators
      end
    end
  end
end