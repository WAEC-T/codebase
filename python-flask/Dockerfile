FROM python:3.10-slim

# Set the working directory
WORKDIR /app

# Install system dependencies for building Python packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libpq-dev \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Copy the application code and requirements
COPY requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application code
COPY . /app

# Expose the application port
EXPOSE 5000

# Set Flask environment variable
ENV FLASK_APP=wsgi.py

# Run the application using Gunicorn
CMD ["gunicorn", "--bind", "0.0.0.0:5000", "wsgi:app"]