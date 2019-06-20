module Provider
  module Azure
    module Config

      # The configuration file path, this should be the root path relative to
      # all API definitions, overrides and examples.
      attr_reader :config_file

      # Azure-extended Provider::Config::validate
      def azure_validate
        @datasources.__is_data_source = true unless @datasources.nil?
      end

      # Azure-extended Provider::Config::parse
      def azure_parse(cfg_file)
        @config_file = cfg_file
      end
    end
  end
end
