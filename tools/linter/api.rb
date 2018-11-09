require 'api/compiler'

class ApiFetcher
  def initialize(filename)
    @filename = filename
    @api = get_yaml
  end

  def fetch
    @api
  end

  private

  def get_yaml
    Api::Compiler.new(@filename).run
  end
end
