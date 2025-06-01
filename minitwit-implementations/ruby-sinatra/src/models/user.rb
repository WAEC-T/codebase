# frozen_string_literal: true

require 'digest'

# Represents a user model in minitwit.
class User < ActiveRecord::Base
  has_many :messages, foreign_key: :author_id

  # manually handle password encryption and authentication
  attr_reader :password

  validates :username, presence: true, uniqueness: true
  validates :email, presence: true, format: { with: URI::MailTo::EMAIL_REGEXP }
  validates :password, presence: true, confirmation: true, if: -> { pw_hash.blank? || !password.nil? }

  # overrides the password= method to set the encrypted password in the pw_hash column
  def password=(unencrypted_password)
    return unless unencrypted_password.present?

    @password = unencrypted_password
    self.pw_hash = unencrypted_password
  end

  # authenticates the user by comparing the provided password with the stored hash
  def authenticate(unencrypted_password)
    pw_hash == unencrypted_password ? self : false
  end

  validates :username, presence: true, uniqueness: true
  validates :password, :password_confirmation, presence: true
  validates :email, presence: true, format: { with: URI::MailTo::EMAIL_REGEXP }

  def follows?(other_user)
    following.include?(other_user)
  end

  def gravatar(size = 80)
    md5 = Digest::MD5.new
    md5 << email.strip.downcase.encode('utf-8')
    md5_hash = md5.hexdigest
    "https://www.gravatar.com/avatar/#{md5_hash}?d=identicon&s=#{size}"
  end
end
