# frozen_string_literal: true

class Follower < ActiveRecord::Base
  belongs_to :who, class_name: "User"
  belongs_to :whom, class_name: "User"
end
