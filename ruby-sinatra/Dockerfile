FROM ruby:3.3.0

RUN apt update && apt upgrade -y

ENV APP_ENV=production
ENV PORT=5000

WORKDIR /waect-ruby-sinatra

COPY Gemfile ./

RUN gem install bundler && bundle install

COPY . ./

EXPOSE 5000

CMD ["bundle", "exec", "ruby", "main.rb"]
