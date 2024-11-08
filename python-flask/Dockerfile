FROM python:3.10-slim

# Set the working directory
WORKDIR /app

# Copy the application code and requirements
COPY . /app
COPY requirements.txt requirements.txt

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Expose the application port
EXPOSE 5000

# Set Flask environment variable
ENV FLASK_APP=wsgi.py

# Run the application using Gunicorn
CMD ["gunicorn", "--bind", "0.0.0.0:5000", "wsgi:app"]