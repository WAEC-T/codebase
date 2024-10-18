FROM python:3.9-slim
WORKDIR /app
COPY . /app
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt
EXPOSE 5000
ENV FLASK_APP=wsgi.py
CMD ["gunicorn", "--bind", "0.0.0.0:5000", "wsgi:app"]
