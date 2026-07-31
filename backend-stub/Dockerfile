FROM python:3.14-alpine

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY app ./app
COPY media ./media
EXPOSE 3001
CMD ["litestar", "--app", "app.main:app", "run", "--host", "0.0.0.0", "--port", "3001"]
