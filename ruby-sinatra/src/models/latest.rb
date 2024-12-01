# frozen_string_literal: true

# Api latest model to store the value in the last operation.
class Latest < ActiveRecord::Base
  self.table_name = 'latest'

  class << self
    def primary_key_name
      'id'
    end

    def value_column_name
      'value'
    end

    def set(latest_value)
      Latest.first.update!(value_column_name => latest_value)
    end

    def get
      Latest.first.send(value_column_name)
    end
  end
end
