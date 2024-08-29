# frozen_string_literal: true

class Latest < ActiveRecord::Base
  self.table_name = "latest"

  class << self
    def set(latest)
      Latest.first.update!(value: latest)
    end

    def get
      Latest.first.value
    end
  end
end
