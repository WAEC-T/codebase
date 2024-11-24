FROM alpine:edge

ENV APP_ENV=production
ENV PORT=5000

RUN apk add --no-cache ruby ruby-dev build-base postgresql-dev && gem install bundler

WORKDIR /waect-ruby-sinatra

COPY Gemfile ./

RUN bundle install

COPY . ./

EXPOSE 5000

CMD ["bundle", "exec", "ruby", "main.rb"]
