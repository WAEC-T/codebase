FROM python:3.9-slim
WORKDIR /app
COPY . /app
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt
EXPOSE 5000
ENV FLASK_APP=wsgi.py
CMD ["flask", "run", "--host=0.0.0.0"]
