FROM alpine:edge
# tesing dockerfile
WORKDIR /app

RUN apk add --no-cache python3 py3-pip gcc musl-dev postgresql-dev

COPY requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

COPY . /app

EXPOSE 5000
ENV FLASK_APP=wsgi.py

CMD ["gunicorn", "--bind", "0.0.0.0:5000", "wsgi:app"]