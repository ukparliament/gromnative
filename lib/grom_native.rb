require 'ffi'
require 'json'
require 'grom'
require 'grom_native/version'
require 'grom_native/node'

# Top level namespace for our gem
module GromNative
  extend FFI::Library
  ffi_lib File.expand_path("../ext/gromnative.so", File.dirname(__FILE__))
  attach_function :get, [:string], :string

  def self.fetch(uri:, headers: {}, filter: [], decorators: nil)
    input = { uri: uri, headers: headers, filter: filter }
    data_struct = JSON.parse(get(uri))

    handle_errors(data_struct)

    build_nodes(data_struct, filter, decorators)
  end

  def self.handle_errors(data_struct)
    error = nil
    status_code = data_struct.fetch('status_code', 0)

    if status_code >= 500
      error = "Server error: #{data_struct['error']}"
    elsif status_code >= 300 && status_code < 500
      error = "Client error: #{data_struct['error']}"
    end

    # require 'irb'; binding.irb
    error ||= data_struct['error']

    raise StandardError, error if error != "" && error != nil
  end

  def self.build_nodes(data_struct, filter, decorators)
    nodes = []
    nodes_by_subject = {}
    filtered_nodes = Array.new(filter.size) { [] }

    data_struct.fetch('statementsBySubject', []).each do |subject, statements|
      node = GromNative::Node.new(statements, decorators)

      nodes << node
      nodes_by_subject[subject] = node

      if !filter.empty? && node.respond_to?(:type)
        node_types = node.blank? ? Array(::Grom::Node::BLANK) : Array(node.type)
        indexes = node_types.reduce([]) do |memo, type|
          index = filter.index(type)
          memo << index if index

          memo
        end

        indexes.each { |index| filtered_nodes[index] << node }
      end
    end

    link_nodes(nodes_by_subject, data_struct)

    return filtered_nodes.first if filter && filter.size == 1
    return filtered_nodes       if filter

    nodes
  end

  def self.link_nodes(nodes_by_subject, data_struct)
    data_struct.fetch('edgesBySubject', {}).each do |subject, predicates|
      predicates.each do |predicate, object_uris|
        raise NamingError if predicate == 'type'

        current_node = nodes_by_subject[subject]
        next if current_node.nil?

        object_uris.each do |object_uri|
          predicate_name_symbol = "@#{predicate}".to_sym

          # Get the current value (if there is one)
          current_value = current_node.instance_variable_get(predicate_name_symbol)

          object = current_value

          # If we have stored a string, and there are objects to link, create an empty array
          current_value_is_string    = current_value.is_a?(String)
          object_is_array_of_strings = object.all? { |entry| entry.is_a?(String) } if object.is_a?(Array)
          object_by_uri              = nodes_by_subject[object_uri]

          object = [] if (current_value_is_string || object_is_array_of_strings) && object_by_uri

          # If the above is correct, and we have an array
          if object.is_a?(Array)
            # Insert the current value (only if this is a new array (prevents possible duplication),
            # the current value is a string, and there are no linked objects to insert)
            object << current_value if object.empty? && current_value_is_string && object_by_uri.nil?

            # Insert linked objects, if there are any
            object << object_by_uri if object_by_uri
          end

          current_node.instance_variable_set(predicate_name_symbol, object)
        end
      end
    end

    self
  end
end
