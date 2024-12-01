# frozen_string_literal: true

# Follower model in minitwit following the database pattern.
class Follower < ActiveRecord::Base
  belongs_to :who, class_name: 'User'
  belongs_to :whom, class_name: 'User'
end
