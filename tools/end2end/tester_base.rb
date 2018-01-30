# Copyright 2017 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

require 'logger'
require 'open3'

# Executes all examples against a real Google Cloud Platform project. The
# account requires to have Owner (or Editor) to all resources being tested.
class TesterBase
  attr_reader :product
  attr_reader :tests
  attr_writer :test_matrix

  ALL_VERIFIER_PHASE_NAME = 'ALL'.freeze
  ALL_VERIFIER_PHASE_RETURNS = {
    'cleanup' => 1,
    'create' => 0,
    'delete' => 1
  }.freeze

  def provider
    raise 'provider must be implemented'
  end

  # rubocop:disable Lint/DuplicateMethods
  def product
    @product.downcase
  end
  # rubocop:enable Lint/DuplicateMethods

  def header
    raise 'header must be implemented'
  end

  def test
    Timers.instance.measure([provider, product]) do
      log_test_start

      raise "No tests on product '#{@product}'" if @tests.nil? || @tests.empty?

      return @tests.reject { |object| skip_if_not_specified(object) }
                   .map { |object| test_object(object) }
                   .all? { |v| v }
    end
  end

  def validate
    @tests.each do |t|
      if t.respond_to?(:validate)
        t.validate
      else
        Logger.instance.log('end2end', 'parser', "#{t['name']} manual test")
      end
    end
  end

  private

  def log_test_start
    log header
    log '=' * 80
    log ''
    log "Testing #{@product} on #{provider.capitalize}"
    log ''
  end

  def skip_if_not_specified(object)
    return false if @test_matrix.empty?

    # If we have the provider or product specified as wildcard we're good
    return false unless @test_matrix.select do |m|
                          m == "#{provider}:" || m == ":#{product}:" \
                            || m == "#{provider}:#{product}:"
                        end.empty?

    # If not we need a match for :product:object
    return false unless @test_matrix.select do |m|
                          m.end_with?(":#{product}:#{object['name']}")
                        end.empty?

    log "Skipping object #{object['name']}"
    true
  end

  def test_object(object)
    log "Testing object #{object['name']}"
    object['phases'].map do |p|
      return false unless phase(p)
      verify_phase(object, p)
    end
    true
  end

  def phase(phase)
    run_and_log_error do
      log "<==== PHASE: #{phase['name']} ====>"

      raise "No tests on phase '#{phase['name']}'" \
        if phase['apply'].nil? || phase['apply'].empty?

      phase['apply'].each do |data|
        run phase['name'], data['run'], data, variables(data['env'] || {})
      end

      log "Tests on phase '#{phase['name']}' successful"
    end
  end

  # A verifier is a command that is run to confirm that the configuration
  # manager actually did what it claimed
  #
  # Format:
  #   verifiers:
  #     - phase: <phase-name>
  #       command: <command>
  #       exits: <exit-codes>
  #
  # All phases verifiers: If you have a command that is the same on all phases,
  # except what it returns based on the phase you can use the special 'ALL' on
  # the phase. If you do that the command will be applied to:
  #  - cleanup: exit=1
  #  - create: exit=0
  #  - delete: exit=1
  def verify_phase(object, phase)
    run_and_log_error do
      verifier = get_verifier(object, phase)
      return true if verifier.nil?

      command = verifier['command']
      expected_exit = expected_exit(phase, verifier)

      execute(
        'verifier', "verifier(#{phase['name']})", { 'exits' => expected_exit },
        {}, [command.split("\n")
                    .map(&:strip)
                    .map { |l| l.gsub('{{run_id}}', Config::RUN_ID.to_s) }
                    .join(' ')]
      )
    end
  end

  def get_verifier(object, phase)
    return nil if object['verifiers'].nil?

    verifiers = object['verifiers'].select do |v|
      [ALL_VERIFIER_PHASE_NAME, phase['name']].include?(v['phase'])
    end

    raise 'Cannot specify "ALL" and phase verifiers at the same time' \
      if verifiers.size > 1

    return nil if verifiers.empty? || skip_all_verifier(verifiers[0], phase)

    verifiers[0]
  end

  def skip_all_verifier(verifier, phase)
    verifier['phase'] == ALL_VERIFIER_PHASE_NAME \
      && !ALL_VERIFIER_PHASE_RETURNS.key?(phase['name'])
  end

  def expected_exit(phase, verifier)
    default_exit = verifier['exits'] || 0
    return default_exit unless ALL_VERIFIER_PHASE_RETURNS.key?(phase['name'])
    return default_exit unless verifier['phase'] == ALL_VERIFIER_PHASE_NAME
    ALL_VERIFIER_PHASE_RETURNS[phase['name']]
  end

  def run_and_log_error
    begin
      yield
    rescue StandardError => e
      color_log Colors::B_RED, e.to_s
      return false
    end

    true
  end

  def command(_data)
    raise 'command has to be implemented'
  end

  def variables(env)
    env.each do |k, v|
      if v.is_a?(String) && v.include?('{{run_id}}')
        env[k] = v.gsub('{{run_id}}', Config::RUN_ID.to_s)
      end
    end

    env
  end

  def run(phase, name, data, variables)
    execute phase, name, data, variables, command(data)
  end

  def execute(phase, name, data, variables, command)
    log_run name, variables, command

    Open3.popen2e(variables, *command) do |_, std_out_and_err, thread|
      output = Array[*std_out_and_err]
      output.each { |line| log "#{phase}(#{name}): #{line}" }

      exit_code = thread.value.exitstatus

      verify_exit_codes (data['exits'] || [0]), exit_code
      verify_outputs data['outputs'], output

      color_log(Colors::B_YELLOW, data['todo']) unless data['todo'].nil?
    end
  end

  def log_run(name, variables, command)
    log '-' * 80
    log "Running: #{name}"
    log "(variables: #{variables.map { |k, v| "#{k} = #{v}" }.join(', ')})" \
      unless variables.empty?
    log command.join(' ')
  end

  def verify_exit_codes(expected_exits, exit_code)
    return if expected_exits.is_a?(Integer) && expected_exits == exit_code
    return if expected_exits.is_a?(Array) && expected_exits.include?(exit_code)
    raise ['Unexpected return code:',
           "expected=#{expected_exits}, actual=#{exit_code}"].join(' ')
  end

  # We need to have at least 1 valid output that matched expected outputs
  # The expected outputs is an array of arrays. The outer array is tested with
  # AND while the inner array is tested with OR.
  #
  # So if:
  #
  #   - - "A"
  #     - "a"
  #   - - "B"
  #     - "b"
  #
  # It requires the output to have both an "A" or "B", upper or lowercase.
  def verify_outputs(expected_outputs, output)
    return if expected_outputs.nil? || expected_outputs.empty?
    expected_outputs.each do |eo_and|
      found = false
      output.each do |o|
        eo_and.each do |eo_or|
          found = output_matches?(eo_or, o)
          break if found
        end

        break if found
      end
      raise "None of expected outputs not found: #{eo_and}" unless found
    end
  end

  def output_matches?(eo_or, o)
    return true if o.include?(eo_or)

    # If the value is a regular expression, try to match it
    return true if Regexp.new([
      '.*',
      eo_or.gsub('[', '\\[').gsub(']', '\\]').gsub('(', '\\(').gsub(')', '\\)'),
      '.*'
    ].join).match?(o)

    false
  end

  def color_log(color, message)
    log [color, ' ', message, ' ', Colors::NC].join
  end

  def log(message)
    Logger.instance.log(provider, product, message)
  end
end
