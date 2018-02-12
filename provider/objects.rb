module Provider
  # A Hash that stores a list of Objects.
  class Objects < Api::Object
    # TODO(alexstephen): Remove consume_api and all logic to replicate
    # Api::HashArray once objects: on YAML files use RubyObject
    def consume_api(api)
      @__api = api
    end

    def validate
      return unless @__objects.nil? # allows idempotency of calling validate
      convert_findings_to_hash
      ensure_keys_are_objects unless @__api.nil?

      unless @objects.nil?
        check_property :objects, Array
        check_property_list :objects, Provider::Resource
      end

      super
    end

    def [](index)
      @__objects[index]
    end

    def each
      return enum_for(:each) unless block_given?
      @__objects.each { |o| yield o }
      self
    end

    def select
      return enum_for(:select) unless block_given?
      @__objects.select { |o| yield o }
      self
    end

    def fetch(key, *args)
      # *args only holds default value. Needs to mimic ::Hash
      if args.empty?
        # KeyErorr will be thrown if key not found
        @__objects&.fetch(key)
      else
        # args[0] will be returned if key not found
        @__objects&.fetch(key, args[0])
      end
    end

    def key?(key)
      @__objects&.key?(key)
    end

    private

    # Converts every variable into @__objects
    def convert_findings_to_hash
      @__objects = {}
      instance_variables.each do |var|
        next if var.id2name.start_with?('@__')
        @__objects[var.id2name[1..-1]] = instance_variable_get(var)
        remove_instance_variable(var)
      end
    end

    def ensure_keys_are_objects
      @__objects.each_key do |type|
        next unless @__api.objects.select { |o| o.name == type }.empty?
        raise [
          "Object #{type} is not a valid type.",
          "Allowed types are: #{@__api.objects.map(&:name)}"
        ].join(' ')
      end
    end
  end

  # Indicates a resource to be described in the <provider>.yaml
  class Resource < Api::Object
    attr_reader :editable
    def validate
      check_optional_property :editable, Boolean
    end
  end

  # A definition of a Resource (for <provider>.yaml) shared between Puppet &
  # Chef.
  class RubyObject < Resource
    attr_reader :create
    attr_reader :update
    attr_reader :delete
    attr_reader :flush
    attr_reader :pre_fetch
    attr_reader :self_link
    attr_reader :collection
    attr_reader :present
    attr_reader :resource_to_request_patch
    attr_reader :return_if_object
    attr_reader :access_api_results
    attr_reader :provider_helpers

    def validate
      check_optional_property :create, String
      check_optional_property :update, String
      check_optional_property :delete, String
      check_optional_property :flush, String
      check_optional_property :access_api_results, Boolean
      check_optional_property :pre_fetch, String
      check_optional_property :self_link, String
      check_optional_property :collection, String
      check_optional_property :present, String
      check_optional_property :resource_to_request_patch, String
      check_optional_property :return_if_object, String
      # TODO(alexstephen): Turn provider_helpers into a proper class instead
      # of a Hash
      check_optional_property :provider_helpers, Hash
      super
    end

    # TODO(alexstephen): Remove when Api::HashArray is fully replaced with
    # Provider::Objects
    def key?(key)
      instance_variables.include? key
    end
  end
end
