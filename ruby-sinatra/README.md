# itu-minitwit

## Ruby setup guide

Requires Ruby 3.3.0, Bundler and Postgres

### Other prerequisites

```bashrc
sudo apt update && sudo apt install libpq-dev
```

### Setup

- Install required gems with `bundle install`
- Initial database with `bundle exec rake db:create` - If the DB doesn't Exists!!!
- Run migrations `bundle exec rake db:migrate` - Not needed!!
- Start the app with `bundle exec ruby myapp.rb`
