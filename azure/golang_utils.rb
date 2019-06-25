require 'google/golang_utils'

module Azure
  module GolangUtils

    def azure_go_literal(value, go_package = nil)
      return "string(#{go_package}#{'.' unless go_package.nil?}#{value})" if value.is_a?(Symbol)
      go_literal value
    end

  end
end
