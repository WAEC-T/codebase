# itu-minitwit

## Ruby setup guide

Requires Ruby 3.3.0 and Bundler

### Other prerequisites

```bashrc
sudo apt update && sudo apt install libpq-dev
```

### Setup

- Install required gems with `bundle install`
- Initial database with `bundle exec rake db:create`
- Run migrations `bundle exec rake db:migrate`
- Start the app with `bundle exec ruby myapp.rb`
